package guts

import (
	"github.com/fitzerc/five-on-four/data"
	"gorm.io/gorm"
)

type LeagueGuts struct {
    db gorm.DB
}

func NewLeagueGuts(db gorm.DB) *LeagueGuts {
    return &LeagueGuts{db: db}
}

func (lg LeagueGuts) Add(newLeague data.League) error {
    return lg.db.Save(&newLeague).Error
}

func (lg LeagueGuts) Delete(id string) error {
    return lg.db.Delete(&data.League{}, id).Error
}

func (lg LeagueGuts) GetById(id string) (data.League, error) {
    var league data.League
    err := lg.db.Where("id = ?", id).First(&league).Error

    return league, err
}

func (lg LeagueGuts) GetAll() ([]data.League, error) {
    var leagues []data.League
    err := lg.db.Find(&leagues).Error

    return leagues, err
}
