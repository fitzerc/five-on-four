package data

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserRole struct {
	gorm.Model
    UserId          uint   `json:"user_id"`
    Role            string `json:"role"`
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
    LeagueName     string `json:"league_name"`
    ActiveSeasonId uint   `json:"active_season_id"`
}
