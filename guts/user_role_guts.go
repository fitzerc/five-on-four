package guts

import (
	"strings"

	"github.com/fitzerc/five-on-four/data"
	"gorm.io/gorm"
)

type UserRoleGuts struct {
	db gorm.DB
}

func NewUserRoleGuts(db gorm.DB) *UserRoleGuts {
	return &UserRoleGuts{db: db}
}

func (urg UserRoleGuts) Save(newUserRole *data.UserRole) error {
	newUserRole.Role = strings.ToLower(newUserRole.Role)

	return urg.db.Save(&newUserRole).Error
}

func (urg UserRoleGuts) GetByQuery(query string, params ...interface{}) ([]data.UserRole, error) {
	var roles []data.UserRole
	err := urg.db.Where(query, params...).Find(&roles).Error

	return roles, err
}

func (urg UserRoleGuts) Delete(id string) error {
	return urg.db.Delete(&data.UserRole{}, id).Error
}
