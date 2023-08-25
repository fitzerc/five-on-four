package handlers

import (
	"net/http"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/labstack/echo/v4"
)

type UserRolesHandler struct {
	UserRoleGuts guts.UserRoleGuts
	UserHandler  UserHandler
}

func (urh UserRolesHandler) RegisterEndpoints(group *echo.Group) {
	group.POST("/users/roles", urh.AddUserRole, urh.UserHandler.MustBeAdmin())
	group.GET("/users/:id/roles", urh.GetRolesByUserId)
	group.DELETE("/users/:id/roles/:roleId", urh.RemoveRoleFromUser, urh.UserHandler.MustBeAdmin())
}

func (roleHandler UserRolesHandler) AddUserRole(c echo.Context) (err error) {
	newRole := new(data.UserRole)
	if err = c.Bind(newRole); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = roleHandler.UserRoleGuts.Save(newRole)
	return c.String(http.StatusOK, "success")
}

func (roleHandler UserRolesHandler) GetRolesByUserId(c echo.Context) (err error) {
	id := c.Param("id")

	roles, err := roleHandler.UserRoleGuts.GetByQuery("user_id = ?", id)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}

	return c.JSON(http.StatusOK, roles)
}

func (roleHandler UserRolesHandler) RemoveRoleFromUser(c echo.Context) (err error) {
	id := c.Param("id")
	roleId := c.Param("roleId")

	_, err = roleHandler.UserRoleGuts.GetByQuery("user_id = ? and id = ?", id, roleId)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}

	err = roleHandler.UserRoleGuts.Delete(roleId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}

	return c.JSON(http.StatusOK, "success")
}
