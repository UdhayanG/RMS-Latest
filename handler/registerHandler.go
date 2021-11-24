package handler

import (
	"RMS-Trail/config"
	"RMS-Trail/domain/form"
	"RMS-Trail/domain/model"
	"crypto/rsa"
	"crypto/tls"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/vonage/vonage-go-sdk"
	"golang.org/x/crypto/bcrypt"
	gomail "gopkg.in/mail.v2"
	"gorm.io/gorm"
)

var (
	cfg config.Properties
)

const (
	tokenExpiresIn = 15000
)
const (
	mySigningKey = "Key,Value"
)

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type TokenHandler struct {
	PrivateKey *rsa.PrivateKey
}

func init() {

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Cannot Read Config : %v", err)
	}
}

// CreateRMS godoc
// @Summary Create a RMS
// @Description Create a new todo item
// @Tags RMS
// @Accept json
// @Produce json
// @Param RMS body form.Registration true "New User"
// @Success 201 {object} form.Registration
// @Router /register [post]
func Register(tx *gorm.DB, verifyClient *vonage.VerifyClient, store *sessions.CookieStore) echo.HandlerFunc {
	return func(c echo.Context) error {
		gob.Register(map[string]interface{}{})
		var u form.Registration
		if err := c.Bind(&u); err != nil {
			return c.JSON(http.StatusUnauthorized, err.Error())
		}
		sess, _ := store.Get(c.Request(), "acmeinc-cookie")
		//sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			Secure:   false,
		}
		u.ErrStatus = "Registration failed"
		sess.Values["Userdetails"] = u
		sess.Save(c.Request(), c.Response())
		tx.Transaction(func(db *gorm.DB) error {
			if err := c.Validate(&u); err != nil {
				u.ErrStatus = err.Error()
				sess.Values["Userdetails"] = u
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusUnauthorized, sess.Values["Userdetails"])
			}
			if strings.Compare(u.Password, u.ConfirmPassword) != 0 {
				u.ErrStatus = "Password Mismatch, Please enter correctly"
				sess.Values["Userdetails"] = u
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusUnauthorized, sess.Values["Userdetails"])
			}

			var login model.Logins
			var loginTypes model.LoginTypes

			db.Select("LoginTypeID").Where("LoginTypeDesc = ?", u.Command).First(&loginTypes)
			var duplicateUserCheck *gorm.DB
			if u.Command == "Email" {
				duplicateUserCheck = db.Debug().Where("UserName = ? ", u.EmailAddress).Find(&login)
			}
			if u.Command == "Phone" {
				duplicateUserCheck = db.Debug().Where("UserName = ? ", u.PhoneNumber).Find(&login)
			}
			if duplicateUserCheck.RowsAffected >= 1 {
				u.ErrStatus = "User Name already exists"
				sess.Values["Userdetails"] = u
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusCreated, sess.Values["Userdetails"])
			}
			if u.Command == "Phone" {
				var phoneNoTypes model.PhoneNoTypes
				var countries model.Countries

				db.Select("PhoneNoTypeID").Where("PhoneNoTypeDesc = ?", "Mobile").First(&phoneNoTypes)
				db.Select("CountryID").Where("CountryName = ?", "INDIA").First(&countries)
				phoneInsert := model.PhoneNumbers{CountryID: countries.CountryID, PhoneNumber: u.PhoneNumber, PhoneNoTypeID: phoneNoTypes.PhoneNoTypeID}
				if err := db.Debug().Save(&phoneInsert).Error; err != nil {
					//fmt.Println(err)
					sess.Values["Userdetails"] = u
					sess.Save(c.Request(), c.Response())
					return c.JSON(http.StatusCreated, sess.Values["Userdetails"])
				}
				fmt.Println(phoneInsert.PhoneNoID)
				userInsert := model.User{FirstName: u.FirstName, MiddleName: u.MiddleName, LastName: u.LastName, DefaultAddressID: 1, DefaultPhoneID: phoneInsert.PhoneNoID}
				if err := db.Debug().Save(&userInsert).Error; err != nil {
					//fmt.Println(err)
					sess.Values["Userdetails"] = u
					sess.Save(c.Request(), c.Response())
					return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
				}
				phoneUserInsert := model.UserPhoneNumbers{UserID: userInsert.UserID, PhoneNoID: phoneInsert.PhoneNoID}
				if err := db.Debug().Save(&phoneUserInsert).Error; err != nil {
					//fmt.Println(err)
					sess.Values["Userdetails"] = u
					sess.Save(c.Request(), c.Response())
					return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
				}
				hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
				if err != nil {
					//fmt.Println(err)
					sess.Values["Userdetails"] = u
					sess.Save(c.Request(), c.Response())
					return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
				}

				loginInsert := model.Logins{UserName: u.PhoneNumber, LoginTypeID: loginTypes.LoginTypeID, UserNameVerified: 0, LoginPasswordSalt: string(hash)}
				if err := db.Debug().Save(&loginInsert).Error; err != nil {
					//fmt.Println(err)
					sess.Values["Userdetails"] = u
					sess.Save(c.Request(), c.Response())
					return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
				}
				loginUserInsert := model.UserLogins{UserID: userInsert.UserID, LoginID: loginInsert.LoginID}
				if err := db.Debug().Save(&loginUserInsert).Error; err != nil {
					//fmt.Println(err)
					sess.Values["Userdetails"] = u
					sess.Save(c.Request(), c.Response())
					return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
				}
				response, errResp, err := verifyClient.Request(u.PhoneNumber, "GoTest", vonage.VerifyOpts{CodeLength: 6, Lg: "en-in", WorkflowID: 4})
				if err != nil {
					fmt.Printf("%#v\n", err)
					//c.Redirect(http.StatusPermanentRedirect, "/phoneverification")
					u.ErrStatus = "registered Successfully, The OTP not sent by registered number please check the number once"
					sess.Values["Userdetails"] = u
					return c.JSON(http.StatusAccepted, sess.Values["Userdetails"])
				} else if response.Status != "0" {
					fmt.Println("Error status " + errResp.Status + ": " + errResp.ErrorText)
					//	c.Redirect(http.StatusPermanentRedirect, "/phoneverification")
					u.ErrStatus = "registered Successfully, The OTP not sent by registered number please check the number once"
					sess.Values["Userdetails"] = u
					return c.JSON(http.StatusAccepted, sess.Values["Userdetails"])
				} else {
					u.RequestID = response.RequestId
					u.ErrStatus = "registered Successfully, The OTP has sent by registered number"
					fmt.Println("Request started: " + response.RequestId)
					sess.Values["Userdetails"] = u
					// redirect to "check" page
					//c.Redirect(http.StatusPermanentRedirect, "/phoneverification")
					return c.JSON(http.StatusAccepted, sess.Values["Userdetails"])
				}
			}
			if u.Command == "Email" {
				emailInsert := model.Emails{EmailAddress: u.EmailAddress}
				if err := db.Debug().Save(&emailInsert).Error; err != nil {
					//fmt.Println(err)
					db.Rollback()
					sess.Values["Userdetails"] = u
					sess.Save(c.Request(), c.Response())
					return c.JSON(http.StatusCreated, sess.Values["Userdetails"])
				}
				fmt.Println(emailInsert.EmailID)
				userInsert := model.User{FirstName: u.FirstName, MiddleName: u.MiddleName, LastName: u.LastName, DefaultAddressID: 1, DefaultPhoneID: 2}
				if err := db.Debug().Save(&userInsert).Error; err != nil {
					//fmt.Println(err)
					db.Rollback()
					sess.Values["Userdetails"] = u
					sess.Save(c.Request(), c.Response())
					return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
				}
				emailUserInsert := model.UserEmails{UserID: userInsert.UserID, EmailID: emailInsert.EmailID}
				if err := db.Debug().Save(&emailUserInsert).Error; err != nil {
					//fmt.Println(err)
					db.Rollback()
					sess.Values["Userdetails"] = u
					sess.Save(c.Request(), c.Response())
					return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
				}
				hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
				if err != nil {
					//fmt.Println(err)
					sess.Values["Userdetails"] = u
					sess.Save(c.Request(), c.Response())
					return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
				}

				loginInsert := model.Logins{UserName: u.EmailAddress, LoginTypeID: loginTypes.LoginTypeID, UserNameVerified: 0, LoginPasswordSalt: string(hash)}
				if err := db.Debug().Save(&loginInsert).Error; err != nil {
					//fmt.Println(err)
					sess.Values["Userdetails"] = u
					sess.Save(c.Request(), c.Response())
					return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
				}
				loginUserInsert := model.UserLogins{UserID: userInsert.UserID, LoginID: loginInsert.LoginID}
				if err := db.Debug().Save(&loginUserInsert).Error; err != nil {
					//fmt.Println(err)
					sess.Values["Userdetails"] = u
					sess.Save(c.Request(), c.Response())
					return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
				}
				sendMail(u, c.Request().Header.Get("Origin"))

			}

			return c.JSON(http.StatusCreated, echo.Map{"status": "Registered Successfully"})

		})
		/*sess.Options.MaxAge = -1
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			return err
		}
		fmt.Println(sess.Values["Userdetails"])
		return c.JSON(http.StatusCreated,"Registered Successfully")*/

		return nil
	}
	//return err
}

func OTPVerify(tx *gorm.DB, verifyClient *vonage.VerifyClient, store *sessions.CookieStore) echo.HandlerFunc {

	return func(c echo.Context) error {
		gob.Register(map[string]interface{}{})
		var u form.Registration
		if err := c.Bind(&u); err != nil {
			return c.JSON(http.StatusUnauthorized, err.Error())
		}
		sess, _ := store.Get(c.Request(), "acmeinc-cookie")
		//sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			Secure:   false,
		}
		sess.Values["Userdetails"] = u
		sess.Save(c.Request(), c.Response())
		fmt.Println(u.RequestID)
		fmt.Println(u.OTP)
		fmt.Println(u.PhoneNumber)
		tx.Transaction(func(db *gorm.DB) error {
			var login model.Logins
			userlogindetails := db.Debug().Where("UserName = ? ", u.PhoneNumber).Find(&login)
			if userlogindetails.RowsAffected >= 1 {
				isAlreadyExists := db.Debug().Where("UserNameVerified = ? ", 1).Find(&login)
				if isAlreadyExists.RowsAffected >= 1 {
					return c.JSON(http.StatusCreated, echo.Map{"res": "Phone Number Already Verified"})
				}
				response, errResp, err := verifyClient.Check(u.RequestID, u.OTP)
				if err != nil {
					fmt.Printf("%#v\n", err)
				} else if response.Status != "0" {
					fmt.Println("Error status " + errResp.Status + ": " + errResp.ErrorText)
				} else {
					fmt.Println("Request complete: " + response.RequestId)
					// Set user as authenticated and return to home page
					sess.Values["registered"] = true
					sess.Save(c.Request(), c.Response())

					login.UserNameVerified = 1
					if err := db.Debug().Save(&login).Error; err != nil {
						//fmt.Println(err)
						db.Rollback()
						return c.JSON(http.StatusInternalServerError, "unauthorised")
					}
					return c.JSON(http.StatusAccepted, "OTP verified Successfully")
				}
			}
			return nil
		})
		return nil
	}

}

// CreateRMS godoc
// @Summary Create a RMS
// @Description Create a new User by Email
// @Tags RMS
// @Accept json
// @Produce json
// @Param RMS body form.Registration true "New User"
// @Success 201 {object} form.Registration
// @Router /registerbyemail [post]
func RegisterWithEmail(tx *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		fmt.Println(c.Request().RequestURI)
		fmt.Println(c.Request().Body)
		fmt.Println(c.Request().Host)
		fmt.Println(c.Request().URL)

		gob.Register(map[string]interface{}{})
		var u form.Registration
		if err := c.Bind(&u); err != nil {
			return c.JSON(http.StatusUnauthorized, err.Error())
		}

		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			Secure:   false,
		}
		fmt.Println(sess.Values["Userdetails"])
		u.ErrStatus = "Registration failed"
		sess.Values["Userdetails"] = u
		sess.Save(c.Request(), c.Response())

		/*tx.Transaction(func(db *gorm.DB) error {
			if err := c.Validate(&u); err != nil {
				u.Err = err.Error()
				sess.Values["Userdetails"] = u
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusUnauthorized, sess.Values["Userdetails"])
			}
			if strings.Compare(u.Password, u.ConfirmPassword) != 0 {
				u.Err = "Password Mismatch, Please enter correctly"
				sess.Values["Userdetails"] = u
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusUnauthorized, sess.Values["Userdetails"])
			}
			var login model.Logins
			duplicateUserCheck := db.Debug().Where("UserName = ? ", u.EmailAddress).Find(&login)
			if duplicateUserCheck.RowsAffected >= 1 {
				u.Err = "User Name already exists"
				sess.Values["Userdetails"] = u
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusCreated, sess.Values["Userdetails"])
			}
			var loginTypes model.LoginTypes
			//var phoneNoTypes model.PhoneNoTypes
			var countries model.Countries
			db.Select("LoginTypeID").Where("LoginTypeDesc = ?", "Email").First(&loginTypes)
			//	db.Select("PhoneNoTypeID").Where("PhoneNoTypeDesc = ?", "Mobile").First(&phoneNoTypes)
			db.Select("CountryID").Where("CountryName = ?", "INDIA").First(&countries)
			emailInsert := model.Emails{EmailAddress: u.EmailAddress}
			if err := db.Debug().Save(&emailInsert).Error; err != nil {
				//fmt.Println(err)
				db.Rollback()
				sess.Values["Userdetails"] = u
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusCreated, sess.Values["Userdetails"])
			}
			fmt.Println(emailInsert.EmailID)
			userInsert := model.User{FirstName: u.FirstName, MiddleName: u.MiddleName, LastName: u.LastName, DefaultAddressID: 1, DefaultPhoneID: 2}
			if err := db.Debug().Save(&userInsert).Error; err != nil {
				//fmt.Println(err)
				db.Rollback()
				sess.Values["Userdetails"] = u
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
			}
			emailUserInsert := model.UserEmails{UserID: userInsert.UserID, EmailID: emailInsert.EmailID}
			if err := db.Debug().Save(&emailUserInsert).Error; err != nil {
				//fmt.Println(err)
				db.Rollback()
				sess.Values["Userdetails"] = u
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
			}
			hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
			if err != nil {
				//fmt.Println(err)
				db.Rollback()
				sess.Values["Userdetails"] = u
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
			}
			loginInsert := model.Logins{UserName: u.EmailAddress, LoginTypeID: loginTypes.LoginTypeID, UserNameVerified: 0, LoginPasswordSalt: string(hash)}
			if err := db.Debug().Save(&loginInsert).Error; err != nil {
				//fmt.Println(err)
				db.Rollback()
				sess.Values["Userdetails"] = u
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
			}
			loginUserInsert := model.UserLogins{UserID: userInsert.UserID, LoginID: loginInsert.LoginID}
			if err := db.Debug().Save(&loginUserInsert).Error; err != nil {
				//fmt.Println(err)
				db.Rollback()
				sess.Values["Userdetails"] = u
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusInternalServerError, sess.Values["Userdetails"])
			}
			fmt.Println("inter")
			sendMail(u, c.Request().Header.Get("Origin"))
			fmt.Println("outer")
			return c.JSON(http.StatusCreated, echo.Map{"status": "Registered Successfully"})

		})*/
		/*sess.Options.MaxAge = -1
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			return err
		}
		fmt.Println(sess.Values["Userdetails"])
		return c.JSON(http.StatusCreated,"Registered Successfully")*/

		return nil
	}
	//return err
}

func RegisterWithSocial(tx *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		var socialuser form.SocialRegistration
		if err := c.Bind(&socialuser); err != nil {
			return c.JSON(http.StatusUnauthorized, err.Error())
		}

		sess, _ := session.Get("session", c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   86400 * 7,
			HttpOnly: true,
			Secure:   false,
		}
		fmt.Println(sess.Values["SocialUserdetails"])
		socialuser.Err = "Registration failed"
		sess.Values["SocialUserdetails"] = socialuser
		sess.Save(c.Request(), c.Response())
		fmt.Println(socialuser.Provider)
		tx.Transaction(func(db *gorm.DB) error {
			var login model.Logins
			duplicateUserCheck := db.Debug().Where("UserName = ? ", socialuser.Id).Find(&login)
			if duplicateUserCheck.RowsAffected >= 1 {
				/*need to handle for future user operation*/
				return c.JSON(http.StatusCreated, echo.Map{"tempRespose": "Already registered user"})
			}
			var loginTypes model.LoginTypes
			//	var countries model.Countries
			//db.Debug().Select("LoginTypeID").Where("LoginTypeDesc = ?", socialuser.LoginType).First(&loginTypes)
			if socialuser.Provider == "GOOGLE" {
				db.Debug().Select("LoginTypeID").Where("LoginTypeDesc = ?", "Social-Google").First(&loginTypes)

			} else if socialuser.Provider == "FACEBOOK" {
				db.Debug().Select("LoginTypeID").Where("LoginTypeDesc = ?", "Social-Facebook").First(&loginTypes)

			} else {
				db.Debug().Select("LoginTypeID").Where("LoginTypeDesc = ?", "Social").First(&loginTypes)

			}

			emailInsert := model.Emails{EmailAddress: socialuser.Email}
			if err := db.Debug().Save(&emailInsert).Error; err != nil {
				return c.JSON(http.StatusCreated, echo.Map{"err": err})
			}
			userInsert := model.User{FirstName: socialuser.FirstName, LastName: socialuser.LastName, DefaultAddressID: 1, DefaultPhoneID: 2}
			if err := db.Debug().Save(&userInsert).Error; err != nil {
				db.Rollback()
				sess.Values["SocialUserdetails"] = socialuser
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusInternalServerError, sess.Values["SocialUserdetails"])
			}
			emailUserInsert := model.UserEmails{UserID: userInsert.UserID, EmailID: emailInsert.EmailID}
			if err := db.Debug().Save(&emailUserInsert).Error; err != nil {
				return c.JSON(http.StatusCreated, echo.Map{"err": err})
			}
			loginInsert := model.Logins{UserName: socialuser.Id, LoginTypeID: loginTypes.LoginTypeID, UserNameVerified: 1}
			if err := db.Debug().Save(&loginInsert).Error; err != nil {
				//fmt.Println(err)
				db.Rollback()
				sess.Values["SocialUserdetails"] = socialuser
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusInternalServerError, sess.Values["SocialUserdetails"])
			}
			loginUserInsert := model.UserLogins{UserID: userInsert.UserID, LoginID: loginInsert.LoginID}
			if err := db.Debug().Save(&loginUserInsert).Error; err != nil {
				//fmt.Println(err)
				db.Rollback()
				sess.Values["SocialUserdetails"] = socialuser
				sess.Save(c.Request(), c.Response())
				return c.JSON(http.StatusInternalServerError, sess.Values["SocialUserdetails"])
			}
			return c.JSON(http.StatusCreated, socialuser)

		})

		return nil

	}

}

//Mail
func sendMail(u form.Registration, ip string) error {

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", cfg.MailFromAddress)

	// Set E-Mail receivers
	m.SetHeader("To", u.EmailAddress)

	// Set E-Mail subject
	m.SetHeader("Subject", "RMS Registration")

	jwttoken, _ := NewToken([]byte(mySigningKey), "RMS Registration", u.EmailAddress)
	mainUrl := ip + "/account/verify-email?token=" + jwttoken
	// Set E-Mail body. You can set plain text or html with text/html
	//m.SetBody("html/template", "<p><a href="+mainUrl+">"+mainUrl+"</a></p>")
	m.SetBody("text/html", fmt.Sprintf("<p>Please click the below link to verify your email address:</p> "+
		"<p><a href= %s>link</a></p>", mainUrl))

	// Settings for SMTP server
	d := gomail.NewDialer(cfg.MailHost, cfg.MailPort, cfg.MailFromAddress, cfg.MailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return nil

}

func NewToken(mySigningKey []byte, subject string, email string) (string, error) {
	createdAt := time.Now().Unix()
	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			Subject:   subject,
			ExpiresAt: createdAt + tokenExpiresIn,
			IssuedAt:  createdAt,
			NotBefore: createdAt,
		},
		Email: email,
	})
	// Set some claims

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString(mySigningKey)

	return tokenString, err
}
func VerifyEmail(tx *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println(c.Request().Header.Get("Authorization"))
		token := c.Request().Header.Get("Authorization")
		fmt.Println(token)
		claims, status, err := Verify(token)

		if err != nil {
			return c.JSON(http.StatusUnauthorized, "invalid request")
		}

		tx.Transaction(func(db *gorm.DB) error {
			var login model.Logins
			userlogindetails := db.Debug().Where("UserName = ? ", claims.Email).Find(&login)
			if userlogindetails.RowsAffected >= 1 {
				isAlreadyExists := db.Debug().Where("UserNameVerified = ? ", 1).Find(&login)
				if isAlreadyExists.RowsAffected >= 1 {
					return c.JSON(http.StatusCreated, echo.Map{"res": "Email Already Verified"})
				}
				login.UserNameVerified = 1
				if err := db.Debug().Save(&login).Error; err != nil {
					//fmt.Println(err)
					db.Rollback()
					return c.JSON(http.StatusInternalServerError, "unauthorised")
				}
				return c.JSON(http.StatusCreated, login)
			} else {
				return c.JSON(http.StatusUnauthorized, echo.Map{"res": "Email Not Registered"})
			}
		})
		fmt.Println(claims.Email)
		fmt.Println(claims.Subject)
		fmt.Println(status)
		return nil
	}
}

func Verify(token string) (*Claims, string, error) {
	fmt.Println(token)
	parsed, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})

	if err != nil {
		return nil, "", err
	}

	if claims, ok := parsed.Claims.(*Claims); ok && parsed.Valid {
		return claims, "valid", nil
	}
	return nil, "", errors.New("invalid token")
}

func SignIn(tx *gorm.DB) echo.HandlerFunc {

	return func(c echo.Context) error {
		var loginForm form.Login
		if err := c.Bind(&loginForm); err != nil {
			return c.JSON(http.StatusUnauthorized, err.Error())
		}
		fmt.Println(loginForm)
		tx.Transaction(func(db *gorm.DB) error {
			var loginModel model.Logins
			usernameList := db.Debug().Where("UserName = ? ", loginForm.UserName).Find(&loginModel)

			if usernameList.RowsAffected == 0 {
				//loginModel.UserName
				fmt.Println("jksdfsjkdfg")
				return c.String(http.StatusUnauthorized, "User Name incorrect")
			}
			if usernameList.RowsAffected >= 1 {
				//encryptpassword, _ := bcrypt.GenerateFromPassword([]byte(loginForm.Password), bcrypt.MinCost)
				//hash, err := bcrypt.GenerateFromPassword([]byte(loginForm.Password), bcrypt.MinCost)
				//fmt.Println(string(encryptpassword))
				err := bcrypt.CompareHashAndPassword([]byte(loginModel.LoginPasswordSalt), []byte(loginForm.Password))
				if err == nil {
					//fmt.Println("sema")
					jwttoken, _ := NewToken([]byte(mySigningKey), "RMS Login", loginForm.UserName)
					fmt.Println(jwttoken)
					return c.JSON(http.StatusOK, echo.Map{"token": jwttoken})

				}
			}
			return nil
		})

		return nil
	}

}

/*func Register(tx *gorm.DB) echo.HandlerFunc{
	return func (c echo.Context) error {
		//To perform a set of Registration operations within a transaction
		tx.Transaction(func(db *gorm.DB) error {
			var u form.Registration
			if err :=c.Bind(&u); err !=nil{
				return c.JSON(http.StatusUnauthorized,err.Error())
			}
			if err := c.Validate(&u); err != nil {
				fmt.Println(err)
				return c.JSON(http.StatusUnauthorized,err.Error())
			}
			if strings.Compare(u.Password, u.ConfirmPassword)!=0{
				return c.JSON(http.StatusUnauthorized,"Password Mismatch, Please enter correctly")

			}
			var login model.Logins
			duplicateUserCheck := db.Debug().Where("UserName = ? ",u.PhoneNumber).Find(&login)
			if(duplicateUserCheck.RowsAffected>=1){
				return c.JSON(http.StatusCreated, "User Name already exists")
			}
			var loginTypes model.LoginTypes
			var phoneNoTypes model.PhoneNoTypes
			var countries model.Countries
				db.Select("LoginTypeID").Where("LoginTypeDesc = ?", "Phone").First(&loginTypes)
				db.Select("PhoneNoTypeID").Where("PhoneNoTypeDesc = ?", "Mobile").First(&phoneNoTypes)
				db.Select("CountryID").Where("CountryName = ?", "INDIA").First(&countries)
				phoneInsert := model.PhoneNumbers{CountryID:countries.CountryID,PhoneNumber:u.PhoneNumber,PhoneNoTypeID:phoneNoTypes.PhoneNoTypeID}
				if err := db.Debug().Save(&phoneInsert).Error; err != nil {
					//fmt.Println(err)
					return c.JSON(http.StatusCreated, err)
				}
				fmt.Println(phoneInsert.PhoneNoID)
				userInsert := model.User{FirstName:u.FirstName,MiddleName:u.MiddleName,LastName:u.LastName,DefaultAddressID:1,DefaultPhoneID:phoneInsert.PhoneNoID}
				if err := db.Debug().Save(&userInsert).Error; err != nil {
					//fmt.Println(err)
					return c.JSON(http.StatusInternalServerError, "Registration failed")
				}
				phoneUserInsert := model.UserPhoneNumbers{UserID:userInsert.UserID,PhoneNoID:phoneInsert.PhoneNoID}
				if err := db.Debug().Save(&phoneUserInsert).Error; err != nil {
					//fmt.Println(err)
					return c.JSON(http.StatusInternalServerError, "Registration failed")
				}
				hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
				if err != nil {
					//fmt.Println(err)
					return c.JSON(http.StatusInternalServerError, "Registration failed")
				}
				loginInsert := model.Logins{UserName:u.PhoneNumber,LoginTypeID:loginTypes.LoginTypeID,UserNameVerified:0,LoginPasswordSalt:string(hash)}
				if err := db.Debug().Save(&loginInsert).Error; err != nil {
					//fmt.Println(err)
					return c.JSON(http.StatusInternalServerError, "Registration failed")
				}
				loginUserInsert := model.UserLogins{UserID:userInsert.UserID,LoginID:loginInsert.LoginID}
				if err := db.Debug().Save(&loginUserInsert).Error; err != nil {
					//fmt.Println(err)
					return c.JSON(http.StatusInternalServerError, "Registration failed")
				}
				return c.JSON(http.StatusCreated, "Registered Successfully")
		})
				return nil
	}
	//return err
}
*/
