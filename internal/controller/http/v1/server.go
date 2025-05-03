package v1

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/flastors/jwt-auth-golang/config"
	"github.com/flastors/jwt-auth-golang/pkg/logging"
)

type Server struct {
	useCases UseCases
	config   config.Config
	logger   *logging.Logger
}

func NewServer(useCases UseCases, config config.Config, logger *logging.Logger) *Server {
	return &Server{
		useCases: useCases,
		config:   config,
		logger:   logger,
	}
}

func (s *Server) Serve() {
	r := NewRouter(s.useCases, s.config.App.Auth, s.logger)
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.config.App.Host, s.config.App.Port))
	if err != nil {
		s.logger.Fatal(err)
	}
	server := &http.Server{
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	s.logger.Info(fmt.Sprintf("Server is listening on %s:%s", s.config.App.Host, s.config.App.Port))
	s.logger.Fatal(server.Serve(listener))
}
