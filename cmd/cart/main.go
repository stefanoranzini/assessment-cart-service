package main

import (
	"github.com/rs/zerolog/log"
	"github.com/stefanoranzini/assessment/cart-service/internal/server"
)

func main() {
	server := server.New(9090)
	err := server.Start()

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
