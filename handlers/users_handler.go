package handlers

import (
	"fmt"
	"net/http"

	"github.com/fitzerc/five-on-four/data"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserHandler struct {
    Db gorm.DB
}

func (userHandler UserHandler) GetLoggedInUser(c echo.Context) (err error) {
    claims, err := GetCustomClaims(c)

    if err != nil {
        return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
            ErrorCode: "invalid_token",
            ErrorDescription: err.Error(),
        })
    }

	var existingUser data.User
	userHandler.Db.Where("id = ?", claims.ID).First(&existingUser)

    if existingUser.ID == 0 {
        return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
            ErrorCode: "invalid_request",
            ErrorDescription: "invalid token",
        })
    }

    return c.JSON(http.StatusOK, existingUser)
}

//TODO: access control - research
//  -claims.UserId must have 'admin' role to add a user
//  -still needs to be available for users to sign up
//    but that could be a /signup that calls this w/ an admin token
func (userHandler UserHandler) AddUser(c echo.Context) (err error) {
	newUser := new(data.User)
	if err = c.Bind(newUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var existingUser data.User
	userHandler.Db.Where("email = ?", newUser.Email).First(&existingUser)

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

    fmt.Printf("%+v\n", newUser)
	newUser.Password = string(hashedPassword)
	userHandler.Db.Save(&newUser)

	return c.String(http.StatusOK, "success")
}
