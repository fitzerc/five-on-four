package handlers

import (
	"net/http"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/fitzerc/five-on-four/utils"
	"github.com/labstack/echo/v4"
)

type LeaguesHandler struct {
    LeagueGuts guts.LeagueGuts
    UserGuts   guts.UserGuts
}

//TODO: GetLeaguesByQueryString?

func (lh LeaguesHandler) GetLeagues(c echo.Context) (err error) {
    leagues, err := lh.LeagueGuts.GetAll()

    if err != nil {
        return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
            ErrorCode: "unknown_error",
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
            ErrorCode: "unknown_error",
            ErrorDescription: err.Error(),
        })
    }

    return c.JSON(http.StatusOK, league)
}

//Admin only
func (lh LeaguesHandler) DeleteLeague(c echo.Context) (err error) {
    claims, err := GetCustomClaims(c)

    if err != nil {
        return err
    }

    userIsAdmin, err := lh.UserGuts.IsAdmin(utils.UintToString(claims.ID))

    if err != nil {
        return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
            ErrorCode: "invalid_token",
            ErrorDescription: err.Error(),
        })
    }

    if userIsAdmin {
        id := c.Param("id")
        err = lh.LeagueGuts.Delete(id)

        if err != nil {
            return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
                ErrorCode: "invalid_token",
                ErrorDescription: err.Error(),
            })
        }

        return c.JSON(http.StatusOK, "success")
    }

    return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
        ErrorCode: "unauthorized",
        ErrorDescription: "action not permitted",
    })
}

//Admin only
func (lh LeaguesHandler) AddLeague(c echo.Context) (err error) {
    claims, err := GetCustomClaims(c)

    if err != nil {
        return err
    }

    userIsAdmin, err := lh.UserGuts.IsAdmin(utils.UintToString(claims.ID))

    if err != nil {
        return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
            ErrorCode: "invalid_token",
            ErrorDescription: err.Error(),
        })
    }

    if userIsAdmin {
        newLeague := new(data.League)

        if err = c.Bind(newLeague); err != nil {
            return echo.NewHTTPError(http.StatusBadRequest, err.Error())
        }

        err = lh.LeagueGuts.Add(*newLeague)

        if err != nil {
            return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
                ErrorCode: "unknown_error",
                ErrorDescription: err.Error(),
            })
        }

        return c.String(http.StatusOK, "success")
    }

    return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
        ErrorCode: "unauthorized",
        ErrorDescription: "action not permitted",
    })
}
