package server

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Server struct {
	port int
	mux  *http.ServeMux
}

func New(port int) *Server {
	serverMux := http.NewServeMux()

	serverMux.HandleFunc("/order", insertOrder)

	return &Server{
		port: port,
		mux:  serverMux,
	}
}

func (s *Server) Start() error {
	log.Info().Msgf("Starting server on port %d", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.mux)
}
