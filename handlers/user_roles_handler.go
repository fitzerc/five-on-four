package handlers

import (
	"net/http"
	"strings"

	"github.com/fitzerc/five-on-four/data"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type UserRolesHandler struct {
    Db gorm.DB
}

//TODO: access control
// -claims.UserId must have 'admin' role
// TODO: update to allow list of roles too
func (roleHandler UserRolesHandler) AddUserRole(c echo.Context) (err error) {
    newRole := new(data.UserRole)
    if err = c.Bind(newRole); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error())
    }

    newRole.Role = strings.ToLower(newRole.Role)

    roleHandler.Db.Save(&newRole)
    return c.String(http.StatusOK, "success")
}

//TODO: access control
// -claims.UserId must have 'admin' role or
//  claims.UserId must equal the id passed in query string
func (roleHandler UserRolesHandler) GetRolesByUserId(c echo.Context) (err error){
    id := c.Param("id")

    var roles []data.UserRole
    err = roleHandler.Db.Where("user_id = ?", id).Find(&roles).Error

    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error)
    }

    return c.JSON(http.StatusOK, roles)
}

//TODO: access control
// -claims.UserId must have 'admin' role
func (roleHandler UserRolesHandler) RemoveRoleFromUser(c echo.Context) (err error){
    id := c.Param("id")
    roleId := c.Param("roleId")

    var role data.UserRole
    err = roleHandler.Db.Where("user_id = ? and id = ?", id, roleId).First(&role).Error

    if err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, err.Error)
    }

    roleHandler.Db.Delete(&data.UserRole{}, roleId)

    return c.JSON(http.StatusOK, "success")
}
