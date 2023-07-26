package data

import (
	"gorm.io/gorm"
)

func InitDb(db *gorm.DB) {
	db.AutoMigrate(&User{}, &UserRole{}, &ReadReceipt{})
}