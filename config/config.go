package config

import (
	"log/slog"

	postgresql "github.com/flastors/jwt-auth-golang/pkg/client/postgres"

	"os"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	configPath = "config/config.yml"
)

type Config struct {
	App     AppConfig     `yaml:"app"`
	Storage StorageConfig `yaml:"storage"`
	SMTP    SMTPConfig    `yaml:"smtp"`
}

type AppConfig struct {
	LogLevel string     `yaml:"log_level" env-default:"debug"`
	Prod     bool       `yaml:"prod" env-default:"false"`
	Host     string     `yaml:"host" env-default:"localhost"`
	Port     string     `yaml:"port" env-default:"8080"`
	Auth     AuthConfig `yaml:"auth"`
}

type AuthConfig struct {
	Secret          string `yaml:"secret"`
	AccessLifetime  int    `yaml:"access_lifetime"`
	RefreshLifetime int    `yaml:"refresh_lifetime"`
}

type StorageConfig struct {
	PostgresConfig postgresql.StorageConfig `yaml:"postgres"`
}

type SMTPConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

var instance *Config
var once sync.Once

func Get() *Config {
	once.Do(func() {
		instance = &Config{}
		readErr := cleanenv.ReadConfig(configPath, instance)
		if readErr != nil {
			description, descrErr := cleanenv.GetDescription(instance, nil)
			if descrErr != nil {
				panic(descrErr)
			}
			slog.Info(description)
			slog.Error(
				"failed to read config",
				slog.String("err", readErr.Error()),
				slog.String("path", configPath),
			)
			os.Exit(1)
		}
	})
	return instance
}
