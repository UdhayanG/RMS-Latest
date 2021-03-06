package main

import (
	//"RMS-Trail/datastore"

	"RMS-Trail/datastore"
	"RMS-Trail/handler"
	"RMS-Trail/utils"
	"log"
	"os"

	"encoding/gob"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gopkg.in/go-playground/validator.v9"

	// "RMS-Trail/domain/form"

	//_ "github.com/swaggo/echo-swagger/example/docs" // docs is generated by Swag CLI, you have to import it.
	//_ "github.com/UdhayanG/Test/docs"
	_ "RMS-Trail/docs"

	"github.com/vonage/vonage-go-sdk"
)

const (
	mySigningKey = "Key,Value"
)

var (
	// Key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)
var verifyClient *vonage.VerifyClient

// @title RMS Application
// @description This is a Repair Management Application
// @version 1.0
// @host localhost:3000
// @BasePath /api/v1
func main() {
	//dsn := "root:Ics@2020@tcp(13.126.170.121:3306)/rms-golang?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Error)})
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pw := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	db, err := datastore.NewDB(host, user, pw, port, dbName)
	if err != nil {
		panic("failed to connect db")
	}
	apiKey := os.Getenv("VONAGE_API_KEY")
	apiSecret := os.Getenv("VONAGE_API_SECRET")

	auth := vonage.CreateAuthFromKeySecret(apiKey, apiSecret)
	verifyClient = vonage.NewVerifyClient(auth)
	//db, err := datastore.NewDB()
	//logFatal(err)
	//db.LogMode(true)
	//db.Logger.LogMode(logger.Info)
	//defer db.Close()
	gob.Register(map[string]interface{}{})
	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	e.Pre(middleware.HTTPSRedirect())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	//	e.AutoTLSManager.HostPolicy = autocert.HostWhitelist("indodev.cf")
	//e.AutoTLSManager.Cache = autocert.DirCache("/var/www/html/GoApp/.cache")
	v1 := e.Group("/api/v1")
	{
		v1.GET("/users", handler.GetUsers(db))
		v1.POST("/check", handler.Check(db))
		v1.POST("/register", handler.Register(db, verifyClient, store))
		v1.POST("/registerbyemail", handler.RegisterWithEmail(db))
		v1.POST("/socialregister", handler.RegisterWithSocial(db))
		v1.POST("/signin", handler.SignIn(db))

	}

	e.POST("/verifyotp", handler.OTPVerify(db, verifyClient, store))
	// Restricted group
	r := e.Group("/api/v1/user/")
	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:     &handler.Claims{},
		SigningKey: []byte(mySigningKey),
	}
	r.Use(middleware.JWTWithConfig(config))

	//r.GET("", handler.GetUserInfo(db))
	e.GET("/account/verify-email", handler.VerifyEmail(db))
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	//err = e.Start(":3000")
	//e.Logger.Fatal(e.StartServer(autotls.DefaultManager("indodevapp.cf").StartAutoTLS(":3000")))
	err = e.StartTLS(":3000", "certificate.crt", "private.key")
	logFatal(err)
}

func logFatal(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
