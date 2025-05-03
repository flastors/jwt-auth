package v1

import (
	"github.com/flastors/jwt-auth-golang/config"
	"github.com/flastors/jwt-auth-golang/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(useCases UseCases, config config.AuthConfig, logger *logging.Logger) *httprouter.Router {
	router := httprouter.New()

	accessHandler := NewAccessHandler(useCases.AccessUseCase, config, logger)
	accessHandler.Register(router)

	refreshHandler := NewRefreshHandler(useCases.RefreshUseCase, config, logger)
	refreshHandler.Register(router)
	return router
}
