package initializers

import (
    "database/sql"
    "fmt"
    "log"

    "github.com/pressly/goose/v3"
)

func MigrateDB(db *sql.DB) {
    // Set the migration directory
    goose.SetBaseFS(nil) // Use the local file system
    goose.SetTableName("goose_db_version") // Set the table name for migration versions

    // Run migrations
    if err := goose.Up(db, "migrations"); err != nil {
        log.Fatalf("goose up failed: %v", err)
    }
    fmt.Println("? Database migrations completed successfully")
}
