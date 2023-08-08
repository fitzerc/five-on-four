package handlers

import (
	"net/http"
	"strconv"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
    UserGuts guts.UserGuts
}

func (userHandler UserHandler) GetLoggedInUser(c echo.Context) (err error) {
    claims, err := GetCustomClaims(c)

    if err != nil {
        return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
            ErrorCode: "invalid_token",
            ErrorDescription: err.Error(),
        })
    }

    //TODO: replace with shared uint to string util
    existingUser, err := userHandler.UserGuts.GetById(
        strconv.FormatUint(uint64(claims.ID), 10))

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

    users, err := userHandler.UserGuts.GetByQuery("email = ?", newUser.Email)

    if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
    }

	if len(users) > 0 && users[0].ID > 0 {
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
    err = userHandler.UserGuts.Save(newUser)

    if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
    }

	return c.String(http.StatusOK, "success")
}
