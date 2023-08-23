package dbutils

import (
	"fmt"
	"github.com/gocraft/dbr/v2"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"path/filepath"
)

func MigrateUp(conn *dbr.Connection, database, migrationSchemaDir string) error {
	if conn == nil {
		panic("nil connection")
	}

	if len(migrationSchemaDir) > 0 {
		tx, err := conn.Begin()

		queryStr := "CREATE DATABASE IF NOT EXISTS " + database
		_, err = tx.Query(queryStr)
		if err != nil {
			return fmt.Errorf("failed to create database '%s': %w", queryStr, err)
		}

		database += "/main"
		_, err = tx.Query("USE " + database)
		if err != nil {
			return fmt.Errorf("failed to use database '%s': %w", database, err)
		}

		tx.Commit()

		driver, err := mysql.WithInstance(conn.DB, &mysql.Config{DatabaseName: database})
		if err != nil {
			return fmt.Errorf("failed to create migration driver: %w", err)
		}

		absPath, _ := filepath.Abs(migrationSchemaDir)
		migrate, err := migrate.NewWithDatabaseInstance("file://"+absPath, database, driver)
		if err != nil {
			return fmt.Errorf("failed to create migration instance: %w", err)
		}

		err = migrate.Up()
		if err != nil {
			return err
		}
	}

	return nil
}
