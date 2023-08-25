package guts

import (
	"github.com/fitzerc/five-on-four/data"
	"gorm.io/gorm"
)

type TeamGuts struct {
	db gorm.DB
}

func NewTeamGuts(db gorm.DB) *TeamGuts {
	return &TeamGuts{db: db}
}

/*
//implement once plaer roles are implemented
func (pg PlayerGuts) IsRole(id string, role string) (bool, error) {
}
*/

func (tg TeamGuts) Add(newTeam data.Team) error {
	return tg.db.Save(&newTeam).Error
}

func (tg TeamGuts) Delete(id string) error {
	return tg.db.Delete(&data.Team{}, id).Error
}

func (tg TeamGuts) GetById(id string) (data.Team, error) {
	var team data.Team
	err := tg.db.Where("id = ?", id).First(&team).Error

	return team, err
}

func (tg TeamGuts) GetByQuery(query string, args ...interface{}) ([]data.Team, error) {
	var teams []data.Team
	err := tg.db.Where(query, args...).Find(&teams).Error

	return teams, err
}

func (tg TeamGuts) GetAll() ([]data.Team, error) {
	var teams []data.Team
	err := tg.db.Find(&teams).Error

	return teams, err
}
