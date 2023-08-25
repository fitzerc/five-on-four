package guts

import (
	"strings"

	"github.com/fitzerc/five-on-four/data"
	"gorm.io/gorm"
)

type PlayerRoleGuts struct {
	db gorm.DB
}

func NewPlayerRoleGuts(db gorm.DB) *PlayerRoleGuts {
	return &PlayerRoleGuts{db: db}
}

func (prg PlayerRoleGuts) Save(newPlayerRole *data.PlayerRole) error {
	newPlayerRole.Role = strings.ToLower(newPlayerRole.Role)

	return prg.db.Save(&newPlayerRole).Error
}

func (prg PlayerRoleGuts) GetByQuery(query string, params ...interface{}) ([]data.PlayerRole, error) {
	var roles []data.PlayerRole
	err := prg.db.Where(query, params...).Find(&roles).Error

	return roles, err
}

func (prg PlayerRoleGuts) Delete(id string) error {
	return prg.db.Delete(&data.PlayerRole{}, id).Error
}
