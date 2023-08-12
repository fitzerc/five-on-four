package handlers

import (
	"net/http"

	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/guts"
	"github.com/labstack/echo/v4"
)

type UserRolesHandler struct {
	UserRoleGuts guts.UserRoleGuts
}

func (urh UserRolesHandler) RegisterEndpoints(group *echo.Group) {
	group.POST("/users/roles", urh.AddUserRole)
	group.GET("/users/:id/roles", urh.GetRolesByUserId)
	group.DELETE("/users/:id/roles/:roleId", urh.RemoveRoleFromUser)
}

// TODO: access control
// -claims.UserId must have 'admin' role
// TODO: update to allow list of roles too
func (roleHandler UserRolesHandler) AddUserRole(c echo.Context) (err error) {
	newRole := new(data.UserRole)
	if err = c.Bind(newRole); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = roleHandler.UserRoleGuts.Save(newRole)
	return c.String(http.StatusOK, "success")
}

// TODO: access control
// -claims.UserId must have 'admin' role or
//
//	claims.UserId must equal the id passed in query string
func (roleHandler UserRolesHandler) GetRolesByUserId(c echo.Context) (err error) {
	id := c.Param("id")

	roles, err := roleHandler.UserRoleGuts.GetByQuery("user_id = ?", id)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error)
	}

	return c.JSON(http.StatusOK, roles)
}

// TODO: access control
// -claims.UserId must have 'admin' role
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
