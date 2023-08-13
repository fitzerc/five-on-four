package guts

import (
	"github.com/fitzerc/five-on-four/data"
	"gorm.io/gorm"
)

type UserGuts struct {
	userRoleGuts UserRoleGuts
	db           gorm.DB
}

func NewUserGuts(userRoleGuts UserRoleGuts, db gorm.DB) *UserGuts {
	return &UserGuts{userRoleGuts: userRoleGuts, db: db}
}

func (ug UserGuts) IsAdmin(id string) (bool, error) {
	roles, err := ug.userRoleGuts.GetByQuery("user_id = ?", id)

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

func (ug UserGuts) GetById(id string) (data.User, error) {
	var user data.User
	err := ug.db.Where("id = ?", id).First(&user).Error

	return user, err
}

func (ug UserGuts) GetByQuery(query string, args ...interface{}) ([]data.User, error) {
	var users []data.User
	err := ug.db.Where(query, args).Find(&users).Error
	return users, err
}

func (ug UserGuts) Save(newUser *data.User) error {
	return ug.db.Save(&newUser).Error
}
