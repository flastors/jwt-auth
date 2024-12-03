package app

import (
	"context"

	"github.com/flastors/jwt-auth-golang/internal/adapters/user_email/empty"
	"github.com/flastors/jwt-auth-golang/internal/adapters/user_repository/postgresql"
	"github.com/flastors/jwt-auth-golang/internal/adapters/user_token/jwt"
	"github.com/flastors/jwt-auth-golang/internal/config"
	v1 "github.com/flastors/jwt-auth-golang/internal/controller/http/v1"
	"github.com/flastors/jwt-auth-golang/internal/core/user"
	postgres "github.com/flastors/jwt-auth-golang/pkg/client/postgres"
	"github.com/flastors/jwt-auth-golang/pkg/logging"
	"github.com/julienschmidt/httprouter"
)

type Context struct {
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) Router() *httprouter.Router {
	return v1.NewRouter(c.UseCases(), c.Config(), c.Logger())
}

func (c *Context) UseCases() *v1.Usecases {
	return &v1.Usecases{
		AccessUseCase:  c.AccessUseCase(),
		RefreshUseCase: c.RefreshUseCase(),
	}
}
func (c *Context) AccessUseCase() *user.AccessUseCase {
	return user.NewAccessUseCase(c.UserRepo(), c.UserToken(), c.Logger())
}

func (c *Context) RefreshUseCase() *user.RefreshUseCase {
	return user.NewRefreshUseCase(c.UserRepo(), c.UserToken(), c.UserEmail(), c.Logger())
}

func (c *Context) UserRepo() user.UserRepository {
	logger := c.Logger()
	return postgresql.NewUserRepository(c.DBClient(logger), logger)
}

func (c *Context) UserToken() user.UserToken {
	return jwt.NewUserToken(c.Config())
}

func (c *Context) UserEmail() user.UserEmail {
	return empty.NewUserEmail()
}

func (c *Context) DBClient(logger *logging.Logger) postgres.Client {
	client, err := postgres.NewClient(context.Background(), 3, c.Config().Storage)
	if err != nil {
		logger.Fatal(err)
	}
	return client
}

func (c *Context) Config() *config.Config {
	return config.GetConfig()
}

func (c *Context) Logger() *logging.Logger {
	return logging.GetLogger()
}
