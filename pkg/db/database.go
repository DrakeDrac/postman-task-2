package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect to database
func ConnectDB(dsn string) error {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error connecting to database")
		return err
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)

	log.Println("Connected to database")
	return nil
}

// Close database connection
func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Println("Error getting DB instance")
		return
	}
	sqlDB.Close()
}

// Migrate database
func MigrateDB(models ...interface{}) error {
	err := DB.AutoMigrate(models...)
	if err != nil {
		log.Println("Migration failed")
	}
	return err
}
