package guts

import (
	"github.com/fitzerc/five-on-four/data"
	"gorm.io/gorm"
)

type TeamMessageBoardGuts struct {
	db gorm.DB
}

func NewTeamMessageBoardGuts(db gorm.DB) *TeamMessageBoardGuts {
	return &TeamMessageBoardGuts{db: db}
}

func (tmbg TeamMessageBoardGuts) Add(newTeamMsgBoard data.TeamMessageBoard) error {
	return tmbg.db.Save(&newTeamMsgBoard).Error
}

func (tmbg TeamMessageBoardGuts) Delete(id string) error {
	return tmbg.db.Delete(&data.TeamMessageBoard{}, id).Error
}

func (tmbg TeamMessageBoardGuts) GetById(id string) (data.TeamMessageBoard, error) {
	var teamMessageBoard data.TeamMessageBoard
	err := tmbg.db.Where("id = ?", id).First(&teamMessageBoard).Error

	return teamMessageBoard, err
}

func (tmbg TeamMessageBoardGuts) GetByQuery(query string, args ...interface{}) ([]data.TeamMessageBoard, error) {
	var teamMsgBoards []data.TeamMessageBoard
	err := tmbg.db.Where(query, args...).Find(&teamMsgBoards).Error

	return teamMsgBoards, err
}

func (tmbg TeamMessageBoardGuts) GetAll() ([]data.TeamMessageBoard, error) {
	var teamMessageBoards []data.TeamMessageBoard
	err := tmbg.db.Find(&teamMessageBoards).Error

	return teamMessageBoards, err
}
