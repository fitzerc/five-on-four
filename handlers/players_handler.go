package handlers

import (
	"net/http"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/labstack/echo/v4"
)

type PlayersHandler struct {
	PlayerGuts  guts.PlayerGuts
	UserHandler UserHandler
}

func (ph PlayersHandler) RegisterEndpoints(group *echo.Group) {
	group.POST("/players", ph.AddPlayer, ph.UserHandler.MustBeAdmin())
	group.DELETE("/players/:id", ph.DeletePlayer, ph.UserHandler.MustBeAdmin())
	group.GET("/players/:id", ph.GetPlayerById)
	group.GET("/players", ph.GetPlayers)
}

func (ph PlayersHandler) AddPlayer(c echo.Context) (err error) {
	newPlayer := new(data.Player)

	if err = c.Bind(newPlayer); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = ph.PlayerGuts.Add(*newPlayer)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	return c.String(http.StatusOK, "success")
}

func (ph PlayersHandler) DeletePlayer(c echo.Context) (err error) {
	id := c.Param("id")
	err = ph.PlayerGuts.Delete(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "invalid_token",
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "success")
}

func (ph PlayersHandler) GetPlayerById(c echo.Context) (err error) {
	id := c.Param("id")
	player, err := ph.PlayerGuts.GetById(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, player)
}

func (ph PlayersHandler) GetPlayers(c echo.Context) (err error) {
	players, err := ph.PlayerGuts.GetAll()

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "uknown_error",
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, players)
}
