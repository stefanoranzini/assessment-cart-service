package main

import (
	"github.com/rs/zerolog/log"
	"github.com/stefanoranzini/assessment/cart-service/internal/server"
)

func main() {
	// In a real case scenario those values will be read from the runtime configuration/environment
	server := server.New(9090, "cart.db")
	err := server.Start()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
