package handlers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type JwtCustomClaims struct {
	Email string `json:"email"`
	ID    uint   `json:"id"`
	jwt.StandardClaims
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenHandler struct {
	UserGuts guts.UserGuts
}

func (tokenHandler TokenHandler) GetApiToken(c echo.Context) (err error) {
	loginReq := new(loginRequest)
	if err = c.Bind(loginReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	users, err := tokenHandler.UserGuts.GetByQuery("email = ?", loginReq.Email)

	if len(users) == 0 || users[0].ID == 0 {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "invalid_credentials",
			ErrorDescription: "unable to login",
		})
	}

	user := users[0]

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password))

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "invalid_credentials",
			ErrorDescription: "unable to login",
		})
	}

	claims := &JwtCustomClaims{
		user.Email,
		user.ID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("SECRET_KEY")
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func GetCustomClaims(c echo.Context) (JwtCustomClaims, error) {
	tokenString := c.Get("user").(*jwt.Token).Raw

	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

	if err != nil {
		return JwtCustomClaims{}, err
	}

	claims, ok := token.Claims.(*JwtCustomClaims)

	if ok {
		return *claims, nil
	}

	return *claims, errors.New("unable to get claims from token")
}

func JWTErrorChecker(err error, c echo.Context) error {
	// Redirects to the signIn form.
	return echo.ErrUnauthorized
}
