package utils

import (
	"errors"
	"github.com/go-playground/locales/id"
	translator "github.com/go-playground/universal-translator"
	id_translations  "gopkg.in/go-playground/validator.v9/translations/en"
	"gopkg.in/go-playground/validator.v9"
	"fmt"
)
// CustomValidator
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate overriding
func (cv *CustomValidator) Validate(i interface{}) error {
	id := id.New()
	uni := translator.New(id, id)

	trans, _ := uni.GetTranslator("id")
	id_translations.RegisterDefaultTranslations(cv.Validator, trans)
	err := cv.Validator.Struct(i)
	
	var errRes string
	fmt.Println(err)
	if err != nil {
		object, _ := err.(validator.ValidationErrors)

		for _, key := range object {
			//fmt.Println(key)
			//errRes +=`
			//,`
			errRes += key.Translate(trans)+","
						
			fmt.Println(errRes)
			//return errors.New(string(key.Translate(trans)))
		}
		return errors.New(errRes)
	}
	/*var b bytes.Buffer
	fmt.Println(err)
	if err != nil {
		object, _ := err.(validator.ValidationErrors)

		for _, key := range object {
			//fmt.Println(key)
			b.WriteString(key.Translate(trans))
			
			//return errors.New(string(key.Translate(trans)))
		}
		return errors.New(b.String())
	}*/
/*	err := cv.Validator.Struct(i)
	if err != nil {

		// translate all error at once
		errs := err.(validator.ValidationErrors)

		// returns a map with key = namespace & value = translated error
		// NOTICE: 2 errors are returned and you'll see something surprising
		// translations are i18n aware!!!!
		// eg. '10 characters' vs '1 character'
		fmt.Println(errs.Translate(trans))
	}*/

	return nil
}
/*package utils

import (
	"unicode"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

// CustomValidator is type setting of third party validator
type CustomValidator struct {
	Validator *validator.Validate
}

// Init validator
func (cv *CustomValidator) Init() {
	cv.Validator.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		var (
			hasNumber      = false
			hasSpecialChar = false
			hasLetter      = false
			hasSuitableLen = false
		)

		password := fl.Field().String()

		if utf8.RuneCountInString(password) <= 30 || utf8.RuneCountInString(password) >= 6 {
			hasSuitableLen = true
		}

		for _, c := range password {
			switch {
			case unicode.IsNumber(c):
				hasNumber = true
			case unicode.IsPunct(c) || unicode.IsSymbol(c):
				hasSpecialChar = true
			case unicode.IsLetter(c) || c == ' ':
				hasLetter = true
			default:
				return false
			}
		}

		return hasNumber && hasSpecialChar && hasLetter && hasSuitableLen
	})
}


    /*if castedObject, ok := err.(validator.ValidationErrors); ok {
        for _, err := range castedObject {
            switch err.Tag() {
            case "required":
                report.Message = fmt.Sprintf("%s is required", 
                    err.Field())
            case "email":
                report.Message = fmt.Sprintf("%s is not valid email", 
                    err.Field())
            case "gte":
                report.Message = fmt.Sprintf("%s value must be greater than %s",
                    err.Field(), err.Param())
            case "lte":
                report.Message = fmt.Sprintf("%s value must be lower than %s",
                    err.Field(), err.Param())
            }

            break
        }
    }*/

// Validate Data
/*func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}*/

