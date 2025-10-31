package main

import (
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"postman-task/internal/api"
	"postman-task/internal/auth"
	"postman-task/internal/core"
	"postman-task/pkg/config"
	"postman-task/pkg/db"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.Load()

	// Connect to database using url from config
	err := db.ConnectDB(cfg.Database.URL)
	if err != nil {
		log.Fatal("error in connecting to database")
	}
	defer db.CloseDB()

	// Load the models and migrate db to use latest schema
	err = db.MigrateDB(&core.User{}, &core.LeaveRequest{}, &core.Attendance{})
	if err != nil {
		log.Println("error in migration")
	}

	// Check if admin user exists
	var adminUser core.User
	err = db.DB.Model(&core.User{}).Where("email = ?", cfg.Admin.Email).First(&adminUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Hash the admin password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.Admin.Password), 10)
			if err != nil {
				log.Fatalf("Failed to hash admin password: %v", err)
			}

			// Create admin user as it was not found
			admin := core.User{
				Name:     "Admin",
				Email:    cfg.Admin.Email,
				Password: string(hashedPassword),
				Role:     "admin",
				Dept:     "Administration",
			}

			// Save admin user to database
			err = db.DB.Create(&admin).Error
			if err != nil {
				log.Fatalf("Failed to create admin user: %v", err)
			}
			log.Printf("Created admin user with email: %s", cfg.Admin.Email)
		} else {
			log.Fatalf("Error checking for admin user: %v", err)
		}
	}

	// Setup jwt auth
	jwt := auth.NewJWTManager(cfg.JWT.SecretKey)

	// Set release mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode("release")
	}

	// Create gin router
	r := gin.Default()

	// Simple health check
	r.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Setup routes
	api.SetupRoutes(r, db.DB, jwt)

	// Start server
	port := "8080"
	if cfg.Server.Port != "" {
		port = cfg.Server.Port
	}

	log.Printf("Starting server on port %s\n", port)
	err = r.Run(":" + port)
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
