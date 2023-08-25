package guts

import (
	"github.com/fitzerc/five-on-four/data"
	"gorm.io/gorm"
)

type PlayerGuts struct {
	db gorm.DB
}

func NewPlayerGuts(db gorm.DB) *PlayerGuts {
	return &PlayerGuts{db: db}
}

func (pg PlayerGuts) Add(newPlayer data.Player) error {
	return pg.db.Save(&newPlayer).Error
}

func (pg PlayerGuts) Delete(id string) error {
	return pg.db.Delete(&data.Player{}, id).Error
}

func (pg PlayerGuts) GetById(id string) (data.Player, error) {
	var player data.Player
	err := pg.db.Where("id = ?", id).First(&player).Error

	return player, err
}

func (pg PlayerGuts) GetByQuery(query string, args ...interface{}) ([]data.Player, error) {
	var players []data.Player
	err := pg.db.Where(query, args...).Find(&players).Error

	return players, err
}

func (pg PlayerGuts) GetAll() ([]data.Player, error) {
	var players []data.Player
	err := pg.db.Find(&players).Error

	return players, err
}
