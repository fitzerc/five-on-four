package data

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string
	Password  string
	FirstName string
	LastName  string
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