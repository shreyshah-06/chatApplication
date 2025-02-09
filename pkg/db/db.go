package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var DB *sql.DB

func InitPostgres() {
	var err error

	// Since we're running from /backend and migrations are in /backend/migrations
	// we should use a relative path from the current working directory
	migrationsPath := "./migrations"  // Changed to point to migrations from backend root
	
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		log.Fatalf("Error getting absolute path to migrations: %v", err)
	}
	fmt.Println("Absolute path to migrations folder:", absPath)

	// Create migration URL using file protocol
	migrationURL := fmt.Sprintf("file://%s", filepath.ToSlash(absPath))
	
	// Print the final migration URL for debugging
	fmt.Println("Final migration URL:", migrationURL)

	// Retrieve the connection string (DSN) from the environment variable
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL not set")
	}

	// Open the database connection
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Check if the database connection is working
	err = DB.Ping()
	if err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	// Initialize the migration instance
	m, err := migrate.New(
		migrationURL,
		dsn,
	)
	if err != nil {
		log.Fatalf("Failed to initialize migrate instance: %v", err)
	}

	// Run the migrations (it will apply any new migrations)
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migrations applied successfully!")
	fmt.Println("Connected to the database")
}