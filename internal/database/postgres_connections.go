package database

import (
	"dms-api/config"
	"dms-api/internal/models"
	"fmt"

	//	"dms-api/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func ConnectionDB() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Config("DB_HOST"), config.Config("DB_USER"), config.Config("DB_PASSWORD"), config.Config("DB_NAME"), config.Config("DB_PORT"))

	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	Database = db
	Database.AutoMigrate(&models.Accounts{}, &models.User{})
	Database.AutoMigrate(&models.OTP{})

	log.Println("Database connected successfully!")
	log.Println("Database migration complete.")
	return nil

}
