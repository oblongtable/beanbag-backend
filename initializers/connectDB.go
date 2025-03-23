package initializers

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
)

// NewDBConnection establishes a new database connection.
func NewDBConnection(config *Config) (*sql.DB, error) {
	// Corrected line: Changed DBname to dbname
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/London", config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)

	DB, err := sql.Open("postgres", dsn) // Use sql.Open
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable the uuid-ossp extension
	if err := enableUUIDExtension(DB); err != nil {
		return nil, fmt.Errorf("failed to enable uuid-ossp extension: %w", err)
	}
	fmt.Println("? uuid-ossp extension enabled")

	fmt.Println("? Connected Successfully to the Database")
	return DB, nil
}

func enableUUIDExtension(DB *sql.DB) error {
	// Check if the extension already exists
	var exists bool
	if err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_extension WHERE extname = 'uuid-ossp')").Scan(&exists); err != nil {
		return err
	}

	if !exists {
		// Enable the extension if it doesn't exist
		if _, err := DB.Exec("CREATE EXTENSION \"uuid-ossp\""); err != nil {
			return err
		}
	}
	return nil
}
