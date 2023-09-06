package data

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string `json:"email" gorm:"not null"`
	Password  string `json:"password" gorm:"not null"`
	FirstName string `json:"first_name" gorm:"not null"`
	LastName  string `json:"last_name" gorm:"not null"`
	Picture   []byte `json:"picture"`
}

type UserRole struct {
	gorm.Model
	UserId          uint   `json:"user_id" gorm:"not null"`
	Role            string `json:"role" gorm:"not null"`
	RoleDescription string `json:"role_description"`
}

type ReadReceipt struct {
	gorm.Model
	UserId  uint `json:"user_id"`
	HasRead bool `gorm:"default:false"`
}

type ErrorResponse struct {
	ErrorCode        string `json:"error_code"`
	ErrorDescription string `json:"error_description"`
}

type League struct {
	gorm.Model
	LeagueName     string `json:"league_name" gorm:"not null"`
	ActiveSeasonId uint   `json:"active_season_id"`
}

type Season struct {
	gorm.Model
	LeagueId   uint      `json:"league_id" gorm:"not null"`
	SeasonName string    `json:"season_name" gorm:"not null"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
}

type Team struct {
	gorm.Model
	SeasonId uint   `json:"season_id" gorm:"not null"`
	TeamName string `json:"team_name" gorm:"not null"`
}

type TeamMessageBoard struct {
	gorm.Model
	TeamId uint `json:"team_id" gorm:"not null"`
}

type Player struct {
	gorm.Model
	UserId       uint `json:"user_id" gorm:"not null"`
	TeamId       uint `json:"team_id" gorm:"not null"`
	Position     string
	JerseyNumber int
}

type PlayerRole struct {
	gorm.Model
	PlayerId        uint   `json:"player_id" gorm:"not null"`
	Role            string `gorm:"not null"`
	RoleDescription string `json:"role_description"`
}
