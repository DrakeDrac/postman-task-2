package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	Admin    AdminConfig
	Email    EmailConfig
}

type DatabaseConfig struct {
	URL string
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	SecretKey string
}

type AdminConfig struct {
	Email    string `mapstructure:"email" yaml:"email"`
	Password string `mapstructure:"password" yaml:"password"`
}

type EmailConfig struct {
	SMTPHost     string `mapstructure:"smtp_host"`
	SMTPPort     string `mapstructure:"smtp_port"`
	SMTPUsername string `mapstructure:"smtp_username"`
	SMTPPassword string `mapstructure:"smtp_password"`
	FromEmail    string `mapstructure:"from_email"`
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Set defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("jwt.secret_key", "mojakey")
	viper.SetDefault("admin.email", "admin@example.com")
	viper.SetDefault("admin.password", "admin123")

	viper.SetDefault("email.smtp_host", "smtp.gmail.com")
	viper.SetDefault("email.smtp_port", "587")
	viper.SetDefault("email.smtp_username", "")
	viper.SetDefault("email.smtp_password", "")
	viper.SetDefault("email.from_email", "")

	viper.BindEnv("database.url", "DATABASE_URL")

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found, using defaults: %v\n", err)
	} else {
		log.Printf("Using config file: %s\n", viper.ConfigFileUsed())
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("Config error:", err)
	}

	if config.Database.URL == "" {
		config.Database.URL = os.Getenv("DATABASE_URL")
	}

	return &config
}
