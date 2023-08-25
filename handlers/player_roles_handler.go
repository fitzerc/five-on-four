package handlers

import (
	"net/http"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/labstack/echo/v4"
)

type PlayerRolesHandler struct {
	PlayerRoleGuts guts.PlayerRoleGuts
	UserHandler    UserHandler
}

func (prh PlayerRolesHandler) RegisterEndpoints(group *echo.Group) {
	group.POST("/players/roles", prh.AddPlayerRole, prh.UserHandler.MustBeAdmin())
	group.GET("/players/:id/roles", prh.GetRolesByPlayerId)
	group.DELETE("/players/:id/roles/:roleId", prh.RemoveRoleFromUser, prh.UserHandler.MustBeAdmin())
}

func (prh PlayerRolesHandler) AddPlayerRole(c echo.Context) (err error) {
	newRole := new(data.PlayerRole)
	if err = c.Bind(newRole); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = prh.PlayerRoleGuts.Save(newRole)
	return c.String(http.StatusOK, "success")
}

func (prh PlayerRolesHandler) GetRolesByPlayerId(c echo.Context) (err error) {
	id := c.Param("id")

	roles, err := prh.PlayerRoleGuts.GetByQuery("player_id = ?", id)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}

	return c.JSON(http.StatusOK, roles)
}

func (prh PlayerRolesHandler) RemoveRoleFromUser(c echo.Context) (err error) {
	id := c.Param("id")
	roleId := c.Param("roleId")

	_, err = prh.PlayerRoleGuts.GetByQuery("player_id = ? and id = ?", id, roleId)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}
	err = prh.PlayerRoleGuts.Delete(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}

	return c.JSON(http.StatusOK, "success")
}
