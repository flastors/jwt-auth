package app

import (
	"context"

	"github.com/flastors/jwt-auth-golang/config"
	"github.com/flastors/jwt-auth-golang/internal/adapters/user_email/gomail"
	"github.com/flastors/jwt-auth-golang/internal/adapters/user_repository/postgresql"
	"github.com/flastors/jwt-auth-golang/internal/adapters/user_token/jwt"
	v1 "github.com/flastors/jwt-auth-golang/internal/controller/http/v1"
	"github.com/flastors/jwt-auth-golang/internal/core/user"
	migration "github.com/flastors/jwt-auth-golang/migrations"
	postgres "github.com/flastors/jwt-auth-golang/pkg/client/postgres"
	"github.com/flastors/jwt-auth-golang/pkg/logging"
)

type Context struct {
}

func NewContext() *Context {
	return &Context{}
}

func (c *Context) Run() {
	c.Migrate()
	c.Server().Serve()
}

func (c *Context) Server() *v1.Server {
	return v1.NewServer(c.UseCases(), *c.Config(), c.Logger())
}

// App Config
func (c *Context) Config() *config.Config {
	return config.Get()
}

// Migrations
func (c *Context) Migrate() {
	migrator, err := migration.NewMigration(c.Config().Storage.PostgresConfig)
	if err != nil {
		panic(err)
	}
	if err := migrator.Up(); err != nil {
		panic(err)
	}
	migrator.Close()
}

// UseCases
func (c *Context) UseCases() v1.UseCases {
	return v1.UseCases{
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

// Adapters
func (c *Context) UserRepo() user.UserRepository {
	logger := c.Logger()
	return postgresql.NewUserRepository(c.DBClient(logger), logger)
}

func (c *Context) UserToken() user.UserToken {
	return jwt.NewUserToken(c.Config().App.Auth)
}

func (c *Context) UserEmail() user.UserEmail {
	return gomail.NewUserEmail(c.Config().SMTP)
}

// pkg
func (c *Context) DBClient(logger *logging.Logger) postgres.Client {
	client, err := postgres.NewClient(context.Background(), 3, c.Config().Storage.PostgresConfig)
	if err != nil {
		logger.Fatal(err)
	}
	return client
}

func (c *Context) Logger() *logging.Logger {
	return logging.GetLogger()
}
