package handlers

import (
    "net/http"

    "github.com/fitzerc/five-on-four/data"
    "github.com/labstack/echo/v4"
    "gorm.io/gorm"
)

type LeaguesHandler struct {
    Db gorm.DB
}

//TODO: GetAllLeagues
//TODO: GetLeagueById
//TODO: GetLeaguesByQueryString?

func (lh LeaguesHandler) GetLeagues(c echo.Context) (err error) {
    var leagues []data.League
    err = lh.Db.Find(&leagues).Error

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

    var league data.League
    err = lh.Db.First(&league, id).Error

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
    userIsAdmin, err := userIsAdmin(c, &lh.Db)

    if err != nil {
        return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
            ErrorCode: "invalid_token",
            ErrorDescription: err.Error(),
        })
    }

    if userIsAdmin {
        id := c.Param("id")
        lh.Db.Delete(&data.League{}, id)

        return c.JSON(http.StatusOK, "success")
    }

    return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
        ErrorCode: "unauthorized",
        ErrorDescription: "action not permitted",
    })
}

//Admin only
func (lh LeaguesHandler) AddLeague(c echo.Context) (err error) {
    userIsAdmin, err := userIsAdmin(c, &lh.Db)

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

        lh.Db.Save(&newLeague)
        return c.String(http.StatusOK, "success")
    }

    return c.JSON(http.StatusBadRequest, &data.ErrorResponse{
        ErrorCode: "unauthorized",
        ErrorDescription: "action not permitted",
    })
}

func userIsAdmin(c echo.Context, db *gorm.DB) (bool, error) {
    claims, err := GetCustomClaims(c)

    if err != nil {
        return false, err
    }

    var roles []data.UserRole
    err = db.Where("id = ?", claims.ID).Find(&roles).Error

    if err != nil {
        return false, err
    }

    isAdmin := false

    for _, r := range roles {
        if r.Role == "admin" {
            isAdmin = true
        }
    }

    return isAdmin, nil
}
