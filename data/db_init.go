package data

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDb(sqliteDbName string) gorm.DB {
	db, err := gorm.Open(sqlite.Open(sqliteDbName), &gorm.Config{})

    if err != nil {
        panic("failed to connect to database")
    }

	db.AutoMigrate(&User{}, &UserRole{}, &ReadReceipt{})
    initData(db)

    return *db;
}

//TODO: replace hard-coded user with a more secure
//      way to create accounts outside of the api
func initData(db *gorm.DB) {
    var existingUser User
    err := db.First(&existingUser).Error

    if err == nil {
        return
    }

    if err.Error() == "record not found" {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)

        if err != nil {
            panic(err)
        }

        db.Create(&User{
            Email: "admin@admin",
            Password: string(hashedPassword),
            FirstName: "Honorable",
            LastName: "Admin",
        })

        db.Create(&UserRole{
            UserId: 1,
            Role: "admin",
            RoleDescription: "THE admin",
        })

        return
    }

    panic(err)
}
