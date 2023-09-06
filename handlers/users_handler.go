package handlers

import (
	"io"
	"net/http"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/fitzerc/five-on-four/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	UserGuts guts.UserGuts
}

func (uh UserHandler) RegisterEndpoints(group *echo.Group) {
	group.POST("/users", uh.AddUser)
	group.GET("/users", uh.GetLoggedInUser)
	group.PUT("/users", uh.UpdateUser)
	group.PUT("/users/image", uh.SaveImage)
}

func (uh UserHandler) MustBeAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := GetCustomClaims(c)

			if err != nil {
				return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
					ErrorCode:        "invalid_token",
					ErrorDescription: err.Error(),
				})
			}

			isAdmin, err := uh.UserGuts.IsAdmin(utils.UintToString(claims.ID))

			if err != nil {
				return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
					ErrorCode:        "internal_error",
					ErrorDescription: err.Error(),
				})
			}

			if !isAdmin {
				return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
					ErrorCode:        "unauthorized",
					ErrorDescription: "You can't do that",
				})
			}

			return next(c)
		}
	}
}

func (userHandler UserHandler) GetLoggedInUser(c echo.Context) (err error) {
	claims, err := GetCustomClaims(c)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "invalid_token",
			ErrorDescription: err.Error(),
		})
	}

	existingUser, err := userHandler.UserGuts.GetById(utils.UintToString(claims.ID))

	if existingUser.ID == 0 {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "invalid_request",
			ErrorDescription: "invalid token",
		})
	}

	return c.JSON(http.StatusOK, existingUser)
}

func (userHandler UserHandler) SignUp(c echo.Context) (err error) {
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

// TODO: check if admin or is user being changed
func (userHandler UserHandler) SaveImage(c echo.Context) (err error) {
	body := c.Request().Body
	defer body.Close()

	claims, err := GetCustomClaims(c)

	user, err := userHandler.UserGuts.GetById(utils.UintToString(claims.ID))

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	if user.ID == 0 {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "user_not_found",
			ErrorDescription: "user not found",
		})
	}

	user.Picture, err = io.ReadAll(body)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "bad_request",
			ErrorDescription: "unable to read image",
		})
	}

	err = userHandler.UserGuts.UpdateImage(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	return c.String(http.StatusOK, "success")
}

// TODO: check if admin or is user being changed
func (userHandler UserHandler) UpdateUser(c echo.Context) (err error) {
	tokenString := c.Get("user").(*jwt.Token).Raw
	claims, err := GetCustomClaims(c)

	newUser := new(data.User)
	if err = c.Bind(newUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user, err := userHandler.UserGuts.GetById(utils.UintToString(claims.ID))

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	if user.ID == 0 {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "user_not_found",
			ErrorDescription: "user not found",
		})
	}

	newUser.ID = user.ID
	newUser.Password = user.Password

	userHandler.UserGuts.Save(newUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	userResponse := LoginResponse{
		Email:     newUser.Email,
		FirstName: newUser.FirstName,
		LastName:  newUser.LastName,
		Token:     tokenString,
	}

	return c.JSON(http.StatusOK, userResponse)
}

// TODO: access control - research
//
//	-claims.UserId must have 'admin' role to add a user
//	-still needs to be available for users to sign up
//	  but that could be a /signup that calls this w/ an admin token
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

func (uh UserHandler) DeleteUser(c echo.Context) (err error) {
	id := c.Param("id")
	err = uh.UserGuts.Delete(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "invalid_token",
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "success")
}
