package v1

import (
	"github.com/flastors/jwt-auth-golang/internal/config"
	"github.com/flastors/jwt-auth-golang/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(usecases *Usecases, cfg *config.Config, logger *logging.Logger) *httprouter.Router {
	router := httprouter.New()

	accessHandler := NewAccessHandler(usecases.AccessUseCase, cfg, logger)
	accessHandler.Register(router)

	refreshHandler := NewRefreshHandler(usecases.RefreshUseCase, cfg, logger)
	refreshHandler.Register(router)
	return router
}
