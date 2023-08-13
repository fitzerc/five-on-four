package handlers

import (
	"net/http"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/labstack/echo/v4"
)

type SeasonsHandler struct {
	SeasonGuts  guts.SeasonGuts
	UserHandler UserHandler
}

func (sh SeasonsHandler) RegisterEndpoints(group *echo.Group) {
	group.POST("/seasons", sh.AddSeason, sh.UserHandler.MustBeAdmin())
	group.DELETE("/seasons/:id", sh.DeleteSeason, sh.UserHandler.MustBeAdmin())
	group.GET("/seasons/:id", sh.GetSeasonById)
	group.GET("/seasons", sh.GetSeasons)
}

func (sh SeasonsHandler) AddSeason(c echo.Context) error {
	newSeason := new(data.Season)

	if err := c.Bind(newSeason); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := sh.SeasonGuts.Add(*newSeason); err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	return c.String(http.StatusOK, "success")
}

func (sh SeasonsHandler) DeleteSeason(c echo.Context) error {
	id := c.Param("id")
	err := sh.SeasonGuts.Delete(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "invalid_token",
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, "success")
}

func (sh SeasonsHandler) GetSeasonById(c echo.Context) error {
	id := c.Param("id")
	season, err := sh.SeasonGuts.GetById(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, season)
}

func (sh SeasonsHandler) GetSeasons(c echo.Context) error {
	seasons, err := sh.SeasonGuts.GetAll()

	if err != nil {
		return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
			ErrorCode:        "unknown_error",
			ErrorDescription: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, seasons)
}
