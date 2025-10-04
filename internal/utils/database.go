package utils

import (
	"fmt"
	"os"

	"github.com/glebarez/sqlite"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase establishes a connection to the database
func ConnectDatabase() {
	var err error

	// Check for environment variable to determine database type
	env := os.Getenv("GALAXYERP_ENV")
	if env == "" {
		env = "dev" // default to development
	}

	switch env {
	case "dev":
		// Use SQLite for development
		zap.L().Info("Using SQLite for development")
		DB, err = gorm.Open(sqlite.Open("galaxyerp.db"), &gorm.Config{})
		if err != nil {
			zap.L().Fatal("Failed to connect to SQLite database", zap.Error(err))
		}
	case "test":
		fallthrough
	case "prod":
		// Use PostgreSQL for test and production
		host := viper.GetString("DB_HOST")
		port := viper.GetInt("DB_PORT")
		user := viper.GetString("DB_USER")
		password := viper.GetString("DB_PASSWORD")
		dbname := viper.GetString("DB_NAME")
		sslmode := viper.GetString("DB_SSLMODE")

		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=Asia/Shanghai",
			host, user, password, dbname, port, sslmode)

		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			zap.L().Fatal("Failed to connect to PostgreSQL database", zap.Error(err))
		}
	default:
		// Default to SQLite for development
		zap.L().Info("Using SQLite for development (default)")
		DB, err = gorm.Open(sqlite.Open("galaxyerp.db"), &gorm.Config{})
		if err != nil {
			zap.L().Fatal("Failed to connect to SQLite database", zap.Error(err))
		}
	}

	zap.L().Info("Database connected successfully!")
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
