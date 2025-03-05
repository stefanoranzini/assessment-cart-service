package server

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/stefanoranzini/assessment/cart-service/internal/dao"
	"github.com/stefanoranzini/assessment/cart-service/internal/dao/db"
	"github.com/stefanoranzini/assessment/cart-service/internal/order"
)

type Server struct {
	port int
	mux  *http.ServeMux
}

func New(port int, dataSourceName string) *Server {
	serverMux := http.NewServeMux()

	db, err := db.ConnectSQLite3(dataSourceName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}

	serverMux.Handle("/order", &inserOrderHandler{orderService: order.NewOrderService(
		dao.NewProductDao(db),
		dao.NewOrderDao(db),
	)})

	return &Server{
		port: port,
		mux:  serverMux,
	}
}

func (s *Server) Start() error {
	log.Info().Msgf("Starting server on port %d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}
