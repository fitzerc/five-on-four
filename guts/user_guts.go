package guts

import (
	"github.com/fitzerc/five-on-four/data"
	"github.com/fitzerc/five-on-four/repository"
	"gorm.io/gorm"
)

type UserGuts struct {
	userRoleGuts UserRoleGuts
	userRepo     repository.UserRepo
	db           gorm.DB
}

func NewUserGuts(userRoleGuts UserRoleGuts, userRepo repository.UserRepo, db gorm.DB) *UserGuts {
	return &UserGuts{userRoleGuts: userRoleGuts, userRepo: userRepo, db: db}
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
	return ug.userRepo.GetById(id)
}

func (ug UserGuts) GetByQuery(query string, args ...interface{}) ([]data.User, error) {
	return ug.userRepo.GetByQuery(query, args)
}

func (ug UserGuts) Save(newUser *data.User) error {
	return ug.userRepo.Save(newUser)
}

func (ug UserGuts) Delete(id string) error {
	//TODO: delete user's roles too
	return ug.userRepo.Delete(id)

}

func (ug UserGuts) UpdateImage(user data.User) error {
	return ug.userRepo.UpdateImage(user)
}
