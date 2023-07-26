package main

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"github.com/fitzerc/five-on-four/data"
)

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect to database")
	}

	db.AutoMigrate(&data.User{}, &data.UserRole{}, &data.ReadReceipt{})
}
