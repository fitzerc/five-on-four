package handlers

import (
	"net/http"

	"github.com/fitzerc/five-on-four/data"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func AddUserHandler(c echo.Context, db *gorm.DB) (err error) {

	newUser := new(data.User)
	if err = c.Bind(newUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var existingUser data.User
	db.Where("email = ?", newUser.Email).First(&existingUser)

	if existingUser.ID > 0 {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "duplicate_user",
			ErrorDescription: "user already exists",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "invalid",
			ErrorDescription: "password",
		})
	}

	newUser.Password = string(hashedPassword)
	db.Save(&newUser)

	return c.String(http.StatusOK, "success")
}
