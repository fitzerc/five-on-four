package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/fitzerc/five-on-four/handlers"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Err loading .env file")
	}

	dbName := os.Getenv("DB_NAME")
	db := data.InitDb(dbName)

	e := echo.New()

	//TODO: research how to manage these dependencies
	// use wire or do manually
	leagueGuts := guts.NewLeagueGuts(db)
	userRoleGuts := guts.NewUserRoleGuts(db)
	userGuts := guts.NewUserGuts(*userRoleGuts, db)

	userHandler := &handlers.UserHandler{UserGuts: *userGuts}
	tokenHandler := &handlers.TokenHandler{UserGuts: *userGuts}
	userRolesHandler := &handlers.UserRolesHandler{
		UserRoleGuts: *userRoleGuts,
		UserHandler:  *userHandler,
	}
	leaguesHandler := &handlers.LeaguesHandler{
		LeagueGuts:  *leagueGuts,
		UserGuts:    *userGuts,
		UserHandler: *userHandler,
	}

	//Unprotected.
	e.POST("/apitoken", tokenHandler.GetApiToken)

	//map endpoints
	apiGroup := e.Group("/api")
	apiGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:                  &jwt.StandardClaims{},
		SigningKey:              []byte(os.Getenv("SECRET_KEY")),
		TokenLookup:             "header:Authorization",
		ErrorHandlerWithContext: handlers.JWTErrorChecker,
	}))

	apiGroup.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	userHandler.RegisterEndpoints(apiGroup)
	userRolesHandler.RegisterEndpoints(apiGroup)
	leaguesHandler.RegisterEndpoints(apiGroup)

	e.Logger.Fatal(e.Start(":1323"))
}
