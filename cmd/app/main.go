package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/flastors/jwt-auth-golang/internal/app"
	"github.com/flastors/jwt-auth-golang/internal/config"
	migration "github.com/flastors/jwt-auth-golang/migrations/postgresql"
	"github.com/flastors/jwt-auth-golang/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

func main() {
	cfg := config.GetConfig()
	logger := logging.GetLogger()
	logger.Info("Starting the application...")
	migrator, err := migration.NewMigration(cfg.Storage)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create migration instance: %v", err))
	}
	err = migrator.Up()
	if err != nil {
		logger.Warn(fmt.Errorf("failed to migrate: %v", err))
	}
	migrator.Close()

	c := app.NewContext()
	router := c.Router()
	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Http.Host, cfg.Http.Port))
	if err != nil {
		logger.Fatal(err)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Info(fmt.Sprintf("Server is listening on %s:%s", cfg.Http.Host, cfg.Http.Port))
	logger.Fatal(server.Serve(listener))
}
