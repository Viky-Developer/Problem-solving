package config

import (
	"Problem-solving/models"
	"errors"
	"fmt"
	"log"
	"os"

	env "github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --------------------------------------------------------------
// AppConfig holds the configuration values for the application.
// --------------------------------------------------------------
type AppConfig struct {
	Port        string
	DatabaseURL string
}

// -----------------------------------------------------------------------------------
// LoadConfig loads environment variables from a .env file and creates an AppConfig.
// -----------------------------------------------------------------------------------
func LoadConfig() (*AppConfig, error) {

	// Load environment variable from a file
	if err := env.Load(".env"); err != nil {
		return nil, errors.New("failed to load env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, errors.New("missing port env")
	}

	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		return nil, errors.New("missing port for database_url")
	}

	return &AppConfig{
		Port:        port,
		DatabaseURL: databaseUrl,
	}, nil
}

// ------------------------------------------------------------
// DbCreate sets up the database connection using the AppConfig
// ------------------------------------------------------------
func DbCreate(config *AppConfig) (*gorm.DB, error) {
	writeDB, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return writeDB, nil
}

// ----------------------------------------------
// CreateTable creates the table and seeds data.
// ----------------------------------------------
func CreateTable(writeDB *gorm.DB) error {

	model := []interface{}{
		&models.KYC{},
	}

	for _, models := range model {
		if err := writeDB.AutoMigrate(models); err != nil {
			return fmt.Errorf("Failed to migrate:%d", err)
		}
	}

	log.Println("Table Created successfully")

	return nil
}
