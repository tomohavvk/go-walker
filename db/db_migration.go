package db

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/tomohavvk/go-walker/config"
)

func PerformMigration(cfg config.DBConfig) error {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	sourceURL := "file://./db/migration"

	fmt.Println("Starting to performing database migrations")

	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No changes. Skip migration")

			return nil
		} else {
			fmt.Println("Error applying migration:", err)
			return err
		}
	}

	sourceErr, databaseError := m.Close()
	if sourceErr != nil {
		fmt.Println("Error closing migration:", sourceErr)
		return sourceErr
	}

	if databaseError != nil {
		fmt.Println("Error closing migration:", databaseError)
		return databaseError
	}

	fmt.Println("Migrations applied successfully.")

	return nil
}
