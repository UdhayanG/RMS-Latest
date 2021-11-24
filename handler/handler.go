package handler

import (
	"RMS-Trail/domain/model"
	"net/http"

	//"github.com/jinzhu/gorm"
	"gorm.io/gorm"
	"github.com/labstack/echo/v4"

)

func Welcome() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome")
	}
}

func GetUsers(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var u []*model.User

		if err := db.Find(&u).Error; err != nil {
			// error handling here
			return err
		}

		return c.JSON(http.StatusOK, u)
	}
}
func Check(db *gorm.DB) echo.HandlerFunc {
	return func (c echo.Context) error {
	u := new(*model.User)
	if err := c.Bind(u); err != nil {
		return err
	}
	result := db.Create(&u)
	return c.JSON(http.StatusCreated, result)
}
}




