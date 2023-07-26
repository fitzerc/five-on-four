package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/fitzerc/five-on-four/data"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type jwtCustomClaims struct {
	Email string `json:"email"`
	ID    uint   `json:"id"`
	jwt.RegisteredClaims
}

type loginRequest struct {
	Email    string `json: "email"`
	Password string `json: "password"`
}

func GetApiTokenHandler(c echo.Context, db *gorm.DB) (err error) {
	loginReq := new(loginRequest)
	if err = c.Bind(loginReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var user data.User
	db.Where("email = ?", loginReq.Email).First(&user)

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

	claims := &jwtCustomClaims{
		user.Email,
		user.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
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
