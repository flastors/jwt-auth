package migration

import (
	"fmt"

	postgres "github.com/flastors/jwt-auth-golang/pkg/client/postgres"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func NewMigration(sc postgres.StorageConfig) (*migrate.Migrate, error) {
	m, err := migrate.New(
		"file://migrations/postgresql",
		fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database),
	)
	if err != nil {
		return nil, err
	}
	return m, nil

}
