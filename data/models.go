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
	UserID          uint
	Role            string
	RoleDescription string
}

type ReadReceipt struct {
	gorm.Model
	UserID  uint
	HasRead bool `gorm:"default:false"`
}

type ErrorResponse struct {
	ErrorCode        string `json:"error_code"`
	ErrorDescription string `json:"error_description"`
}
