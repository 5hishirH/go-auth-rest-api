package db

import (
	"database/sql"
	"fmt"

	"github.com/5hishirH/go-auth-rest-api.git/migrations"
	_ "github.com/mattn/go-sqlite3"
)

func NewSqliteStorage(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Run Migrations
	if err := applyMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func applyMigrations(db *sql.DB) error {
	// Read all files in the "migrations" folder
	files, err := migrations.FS.ReadDir(".")
	if err != nil {
		return err
	}

	for _, file := range files {
		fmt.Println("Running migration:", file.Name())

		// Read the file content
		content, err := migrations.FS.ReadFile(file.Name())
		if err != nil {
			return err
		}

		// Execute the SQL
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("migration %s failed: %w", file.Name(), err)
		}
	}
	return nil
}
