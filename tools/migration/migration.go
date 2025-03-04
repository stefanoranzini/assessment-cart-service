package main

import (
	"github.com/rs/zerolog/log"
	"github.com/stefanoranzini/assessment/cart-service/internal/dao/db"
	"github.com/stefanoranzini/assessment/cart-service/tools/migration/migrate"
)

func main() {
	sqlLiteDb, err := db.ConnectSQLite3("cart.db")
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to database")
	}
	defer sqlLiteDb.Close()

	err = migrate.MigrateDB(sqlLiteDb)
	if err != nil {
		log.Fatal().Err(err).Msg("unable to migrate database")
	}
}
