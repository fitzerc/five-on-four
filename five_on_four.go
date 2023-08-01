package main

import (
	"log"
	"net/http"
	"os"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/handlers"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Err loading .env file")
	}

	dbName := os.Getenv("DB_NAME")
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	data.InitDb(db)

	e := echo.New()

    //map endpoints
    apiGroup := e.Group("/api")

    apiGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
        Claims: &jwt.StandardClaims{},
        SigningKey: []byte(os.Getenv("SECRET_KEY")),
        TokenLookup: "header:Authorization",
        ErrorHandlerWithContext: handlers.JWTErrorChecker,
    }))

	apiGroup.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.POST("/users", func(c echo.Context) error {
		return handlers.AddUserHandler(c, db)
	})
	e.POST("/apitoken", func(c echo.Context) error {
		return handlers.GetApiTokenHandler(c, db)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
