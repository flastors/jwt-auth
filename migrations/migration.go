package migration

import (
	"fmt"
	"time"

	postgresql "github.com/flastors/jwt-auth-golang/pkg/client/postgres"
	"github.com/flastors/jwt-auth-golang/pkg/utils"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Migrator interface {
	Up() error
	Down() error
	Close() (error, error)
}

func NewMigration(sc postgresql.StorageConfig) (Migrator, error) {
	var m *migrate.Migrate
	err := utils.DoWithRetry(func() error {
		var err error
		m, err = migrate.New(
			"file://migrations",
			fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database),
		)
		return err
	}, 3, 5*time.Second)

	if err != nil {
		return nil, err
	}
	return m, nil

}
