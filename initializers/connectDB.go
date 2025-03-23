package initializers

import (
	"fmt"
	"log"

	"github.com/oblongtable/beanbag-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func ConnectDB(config *Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/London", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // Model 'User' -> Table 'user' rather than table 'users', etc.
		},
	})
	if err != nil {
		log.Fatal("Failed to connect to the Database")
	}
	fmt.Println("? Connected Successfully to the Database")

	// Enable the uuid-ossp extension
	if err := enableUUIDExtension(DB); err != nil {
		log.Fatalf("Failed to enable uuid-ossp extension: %v", err)
	}
	fmt.Println("? uuid-ossp extension enabled")

	// AutoMigrate models
	if err := autoMigrateModels(DB); err != nil {
		log.Fatalf("Failed to migrate models: %v", err)
	}
	fmt.Println("? Models migrated successfully")
}

func enableUUIDExtension(db *gorm.DB) error {
	// Check if the extension already exists
	var exists bool
	if err := db.Raw("SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp')").Scan(&exists).Error; err != nil {
		return err
	}

	if !exists {
		// Enable the extension if it doesn't exist
		if err := db.Exec("CREATE EXTENSION \"uuid-ossp\"").Error; err != nil {
			return err
		}
	}
	return nil
}

func autoMigrateModels(db *gorm.DB) error {
	// List of models to migrate
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Quiz{},
		// Add other models here, e.g., &models.Product{}, &models.Order{}
	}

	// AutoMigrate each model
	for _, model := range modelsToMigrate {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}

	return nil
}
