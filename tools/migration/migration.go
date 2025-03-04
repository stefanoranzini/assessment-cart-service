package main

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"github.com/stefanoranzini/assessment/cart-service/internal/dao/db"
)

//go:embed migrations
var migrations embed.FS

func main() {
	sqlListeDb, err := db.ConnectSQLite3("cart.db")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to database")
	}
	defer sqlListeDb.Close()

	err = migrateDB(sqlListeDb)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to migrate database")
	}
}

func migrateDB(db *sql.DB) error {
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("unable to prepare driver for migration: %v", err)
	}

	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		return fmt.Errorf("unable to prepare source for migration: %v", err)
	}
	migration, err := migrate.NewWithInstance("iofs", source, "sqlite3", driver)
	if err != nil {
		return fmt.Errorf("unable to prepare migration: %v", err)
	}

	err = migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("unable to apply migration: %v", err)
	}

	return nil
}
