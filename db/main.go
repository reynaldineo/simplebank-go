package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Define flags for commands
	action := flag.String("action", "", "Migration action: up, down, up1, down1, create")
	name := flag.String("name", "", "Migration name (required for create action)")
	flag.Parse()

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Database URL from environment variable
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	// Migration files path
	migrationPath := "file://db/migration"

	// Initialize the migrate instance (except for 'create' action)
	var m *migrate.Migrate
	var err error
	if *action != "create" {
		m, err = migrate.New(migrationPath, dbURL)
		if err != nil {
			log.Fatalf("Failed to initialize migrate instance: %v", err)
		}
		defer func() {
			sourceErr, dbErr := m.Close()
			if sourceErr != nil || dbErr != nil {
				log.Fatalf("Error closing migrate instance: sourceErr=%v, dbErr=%v", sourceErr, dbErr)
			}
		}()
	}

	// Perform the requested action
	switch *action {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		fmt.Println("Migrations applied successfully.")

	case "up1":
		if err := m.Steps(1); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply one migration step: %v", err)
		}
		fmt.Println("One migration step applied successfully.")

	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		fmt.Println("Migrations rolled back successfully.")

	case "down1":
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback one migration step: %v", err)
		}
		fmt.Println("One migration step rolled back successfully.")

	case "create":
		// Create a new migration file
		if *name == "" {
			log.Fatal("Migration name is required for the create action")
		}
		createMigration(*name)
		fmt.Println("Migration files created successfully.")

	default:
		log.Fatalf("Invalid action: %v. Valid actions are up, down, up1, down1, create.", *action)
	}
}

func createMigration(name string) {
	// Directory where migrations are stored
	migrationDir := "db/migration"

	// Ensure the migration directory exists
	if err := os.MkdirAll(migrationDir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create migration directory: %v", err)
	}

	// Generate sequential migration IDs
	nextID := getNextMigrationID(migrationDir)

	// Create file names for the .up.sql and .down.sql files
	upFile := fmt.Sprintf("%s/%s_%s.up.sql", migrationDir, nextID, name)
	downFile := fmt.Sprintf("%s/%s_%s.down.sql", migrationDir, nextID, name)

	// Create empty files
	if err := os.WriteFile(upFile, []byte("-- Write your UP migration SQL here\n"), 0644); err != nil {
		log.Fatalf("Failed to create up migration file: %v", err)
	}
	if err := os.WriteFile(downFile, []byte("-- Write your DOWN migration SQL here\n"), 0644); err != nil {
		log.Fatalf("Failed to create down migration file: %v", err)
	}

	fmt.Printf("Migration files created:\n- %s\n- %s\n", upFile, downFile)
}

// Helper to determine the next migration ID
func getNextMigrationID(migrationDir string) string {
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		log.Fatalf("Failed to read migration directory: %v", err)
	}

	// Find the highest sequence number
	highest := 0
	for _, file := range files {
		var seq int
		_, err := fmt.Sscanf(file.Name(), "%d_", &seq)
		if err == nil && seq > highest {
			highest = seq
		}
	}

	return fmt.Sprintf("%06d", highest+1) // Zero-padded 6-digit ID
}
