package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/fitzerc/five-on-four/handlers"
	"github.com/fitzerc/five-on-four/repository"
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
	userRepo := repository.NewUserRepo(db)
	leagueGuts := guts.NewLeagueGuts(db)
	userRoleGuts := guts.NewUserRoleGuts(db)
	userGuts := guts.NewUserGuts(*userRoleGuts, userRepo, db)
	seasonGuts := guts.NewSeasonGuts(db)
	teamGuts := guts.NewTeamGuts(db)
	teamMessageBoardGuts := guts.NewTeamMessageBoardGuts(db)
	playerGuts := guts.NewPlayerGuts(db)
	playerRoleGuts := guts.NewPlayerRoleGuts(db)

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
	seasonsHandler := &handlers.SeasonsHandler{
		SeasonGuts:  *seasonGuts,
		UserHandler: *userHandler,
	}
	teamsHandler := &handlers.TeamsHandler{
		TeamGuts:             *teamGuts,
		TeamMessageBoardGuts: *teamMessageBoardGuts,
		UserHandler:          *userHandler,
	}
	teamMessageBoardHandler := &handlers.TeamMessageBoardsHandler{
		TeamMessageBoardGuts: *teamMessageBoardGuts,
	}
	playersHandler := &handlers.PlayersHandler{
		PlayerGuts:  *playerGuts,
		UserHandler: *userHandler,
	}
	playerRolesHandler := &handlers.PlayerRolesHandler{
		PlayerRoleGuts: *playerRoleGuts,
		UserHandler:    *userHandler,
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
	seasonsHandler.RegisterEndpoints(apiGroup)
	teamsHandler.RegisterEndpoints(apiGroup)
	teamMessageBoardHandler.RegisterEndpoints(apiGroup)
	playersHandler.RegisterEndpoints(apiGroup)
	playerRolesHandler.RegisterEndpoints(apiGroup)

	e.Logger.Fatal(e.Start(":" + os.Getenv("API_PORT")))
}
