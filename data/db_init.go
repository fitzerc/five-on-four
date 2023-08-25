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

	runMigrations(db)
	initData(db)

	return *db
}

func runMigrations(db *gorm.DB) {
	db.AutoMigrate(
		&User{},
		&UserRole{},
		&ReadReceipt{},
		&League{},
		&Season{},
		&Team{},
		&TeamMessageBoard{},
		&Player{},
		&PlayerRole{})
}

// TODO: replace hard-coded user with a more secure
//
//	way to create accounts outside of the api
func initData(db *gorm.DB) error {
	var existingUser User
	err := db.First(&existingUser).Error

	if err != nil && err.Error() == "record not found" {
		err = addDefaultUser(db)
	}

	return err
}

func addDefaultUser(db *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	db.Create(&User{
		Email:     "admin@admin",
		Password:  string(hashedPassword),
		FirstName: "Honorable",
		LastName:  "Admin",
	})

	db.Create(&UserRole{
		UserId:          1,
		Role:            "admin",
		RoleDescription: "THE admin",
	})

	return nil
}
