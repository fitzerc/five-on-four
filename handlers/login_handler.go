package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fitzerc/five-on-four/data"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
    Db gorm.DB
}

func (tokenHandler TokenHandler) GetApiToken(c echo.Context) (err error) {
	loginReq := new(loginRequest)
	if err = c.Bind(loginReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var user data.User
	tokenHandler.Db.Where("email = ?", loginReq.Email).First(&user)

	if user.ID == 0 {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "invalid_credentials",
			ErrorDescription: "unable to login",
		})
	}

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
    fmt.Println("Error:")
    fmt.Printf("%+v", err)
    fmt.Println("Context:")
    fmt.Printf("%+v", c)
	return c.Redirect(http.StatusMovedPermanently, c.Echo().Reverse("userSignInForm"))
}
