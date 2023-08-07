package guts

import (
	"github.com/fitzerc/five-on-four/data"
	"gorm.io/gorm"
)

type UserGuts struct {
    db gorm.DB
}

func NewUserGuts(db gorm.DB) *UserGuts{
    return &UserGuts{db: db}
}

func (ug UserGuts) IsAdmin(id string) (bool, error) {
    var roles []data.UserRole
    err := ug.db.Where("id = ?", id).Find(&roles).Error

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
