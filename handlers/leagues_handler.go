package handlers

import (
	"net/http"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/labstack/echo/v4"
)

type LeaguesHandler struct {
	LeagueGuts  guts.LeagueGuts
	UserGuts    guts.UserGuts
	UserHandler UserHandler
}

func (lh LeaguesHandler) RegisterEndpoints(group *echo.Group) {
	group.POST("/leagues", lh.AddLeague, lh.UserHandler.MustBeAdmin())
	group.DELETE("/leagues/:id", lh.DeleteLeague, lh.UserHandler.MustBeAdmin())
	group.GET("/leagues/:id", lh.GetLeagueById)
	group.GET("/leagues", lh.GetLeagues)
}

// TODO: GetLeaguesByQueryString?
func (lh LeaguesHandler) GetLeagues(c echo.Context) (err error) {
	leagues, err := lh.LeagueGuts.GetAll()

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, leagues)
}

func (lh LeaguesHandler) GetLeagueById(c echo.Context) (err error) {
	id := c.Param("id")
	league, err := lh.LeagueGuts.GetById(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, league)
}

func (lh LeaguesHandler) DeleteLeague(c echo.Context) (err error) {
	id := c.Param("id")
	err = lh.LeagueGuts.Delete(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "invalid_token",
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "success")
}

func (lh LeaguesHandler) AddLeague(c echo.Context) (err error) {
	newLeague := new(data.League)

	if err = c.Bind(newLeague); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = lh.LeagueGuts.Add(*newLeague)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	return c.String(http.StatusOK, "success")
}
