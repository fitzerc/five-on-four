package handlers

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/fitzerc/five-on-four/utils"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginResponse struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Token     string `json:"token"`
}

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

func (tokenHandler TokenHandler) RefreshToken(c echo.Context) (err error) {
	claims, err := GetRefreshCustomClaims(c)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "invalid_token",
			ErrorDescription: err.Error(),
		})
	}

	user, err := tokenHandler.UserGuts.GetById(utils.UintToString(claims.ID))

	t, rt, err := generateTokens(user)

	resp := LoginResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Token:     t,
	}

	c = addRefreshTokenCookie(c, rt)
	c = addAccessTokenCookie(c, t)

	return c.JSON(http.StatusOK, resp)
}

func (th TokenHandler) Logout(c echo.Context) (err error) {
	c = addRefreshTokenCookie(c, "invalid")
	c = addAccessTokenCookie(c, "invalid")

	resp := LoginResponse{
		Token: "",
	}

	return c.JSON(http.StatusOK, resp)
}

func (tokenHandler TokenHandler) Login(c echo.Context) (err error) {
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

	t, rt, err := generateTokens(user)

	retVal := LoginResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Token:     t,
	}

	c = addRefreshTokenCookie(c, rt)
	c = addAccessTokenCookie(c, t)

	return c.JSON(http.StatusOK, retVal)
}

func generateTokens(user data.User) (access string, refresh string, err error) {
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
		return "", "", err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshSecret := os.Getenv("REFRESH_SECRET_KEY")
	rt, err := refreshToken.SignedString([]byte(refreshSecret))

	if err != nil {
		return "", "", err
	}

	return t, rt, nil
}

func addAccessTokenCookie(c echo.Context, token string) echo.Context {
	cookie := new(http.Cookie)
	cookie.Name = "access_token"
	cookie.Value = token
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
	return c
}

func addRefreshTokenCookie(c echo.Context, refreshToken string) echo.Context {
	cookie := new(http.Cookie)
	cookie.Name = "refresh_token"
	cookie.Value = refreshToken
	cookie.Path = "/refresh"
	cookie.Expires = time.Now().Add(30 * (24 * time.Hour))
	c.SetCookie(cookie)
	return c
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

func GetRefreshCustomClaims(c echo.Context) (JwtCustomClaims, error) {
	cookie, err := c.Cookie("refresh_token")

	if err != nil {
		return JwtCustomClaims{}, err
	}

	return getClaimsFromToken(cookie.Value, os.Getenv("REFRESH_SECRET_KEY"))
}

func GetCustomClaims(c echo.Context) (JwtCustomClaims, error) {
	tokenString := c.Get("user").(*jwt.Token).Raw
	return getClaimsFromToken(tokenString, os.Getenv("SECRET_KEY"))
}

func getClaimsFromToken(token_str string, secret string) (JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(token_str, &JwtCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
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
