package guts

import (
	"github.com/fitzerc/five-on-four/data"
	"gorm.io/gorm"
)

type SeasonGuts struct {
	db gorm.DB
}

func NewSeasonGuts(db gorm.DB) *SeasonGuts {
	return &SeasonGuts{db: db}
}

func (sg SeasonGuts) Add(newSeason data.Season) error {
	return sg.db.Save(&newSeason).Error
}

func (sg SeasonGuts) Delete(id string) error {
	return sg.db.Delete(&data.Season{}, id).Error
}

func (sg SeasonGuts) GetById(id string) (data.Season, error) {
	var season data.Season
	err := sg.db.Where("id = ?", id).First(&season).Error

	return season, err
}

func (sg SeasonGuts) GetAll() ([]data.Season, error) {
	var seasons []data.Season
	err := sg.db.Find(&seasons).Error

	return seasons, err
}
