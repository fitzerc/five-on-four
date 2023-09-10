package repository

import (
	"github.com/fitzerc/five-on-four/data"
	"gorm.io/gorm"
)

type UserRepo interface {
	GetById(id string) (data.User, error)
	GetByQuery(query string, args ...interface{}) ([]data.User, error)
	Save(newUser *data.User) error
	Delete(id string) error
	UpdateImage(user data.User) error
}

type GormUserRepo struct {
	db gorm.DB
}

func NewUserRepo(db gorm.DB) UserRepo {
	return &GormUserRepo{db: db}
}

func (gur GormUserRepo) GetById(id string) (data.User, error) {
	var user data.User
	err := gur.db.Where("id = ?", id).First(&user).Error

	return user, err
}

func (gur GormUserRepo) GetByQuery(query string, args ...interface{}) ([]data.User, error) {

	var users []data.User
	err := gur.db.Where(query, args...).Find(&users).Error
	return users, err
}

func (gur GormUserRepo) Save(newUser *data.User) error {
	return gur.db.Save(&newUser).Error
}

func (gur GormUserRepo) Delete(id string) error {
	return gur.db.Delete(&data.User{}, id).Error
}

func (gur GormUserRepo) UpdateImage(user data.User) error {
	return gur.db.Model(&user).Update("picture", []byte(user.Picture)).Error
}
