package data

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
)

func InitDb(sqliteDbName string) gorm.DB {
	db, err := gorm.Open(sqlite.Open(sqliteDbName), &gorm.Config{})

    if err != nil {
        panic("failed to connect to database")
    }

	db.AutoMigrate(&User{}, &UserRole{}, &ReadReceipt{})

    return *db;
}
