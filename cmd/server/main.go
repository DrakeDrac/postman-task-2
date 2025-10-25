package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load configurations from environment variables
	dbUrl := os.Getenv("DATABASE_URL")

	var db *gorm.DB
	var err error
	db, err = gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err == nil {
		fmt.Println("Connected to the database!")
	} else {
		fmt.Println("Failed to connect to database.", err)
	}

	fmt.Println("✅ Database connection established successfully!")

	// init gin
	r := gin.Default()

	r.GET("/test", func(c *gin.Context) {
		// Check if the DB connection is still active
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "message": "Failed to access database"})
			return
		}

		err = sqlDB.Ping()
		if err != nil {
			c.JSON(500, gin.H{"status": "error", "message": "Database ping failed"})
			return
		}

		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Successfully connected to the database",
		})
	})

	fmt.Println("Server started on port 8080...")
	eerr := r.Run(":8080")
	if eerr != nil {
		fmt.Println("Failed to run server: ", eerr)
	}
}
