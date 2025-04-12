package initializers

import (
    "database/sql"
    "fmt"
    "log"
    "embed"

    "github.com/pressly/goose/v3"
)

func MigrateDB(db *sql.DB, migrationsEmbedRoot embed.FS) {
    // Set the migration directory
    goose.SetBaseFS(migrationsEmbedRoot) // Use the embedded dir for goose
    goose.SetTableName("goose_db_version") // Set the table name for migration versions

    // Run migrations
    if err := goose.Up(db, "migrations"); err != nil {
        log.Fatalf("goose up failed: %v", err)
    }
    fmt.Println("? Database migrations completed successfully")
}
