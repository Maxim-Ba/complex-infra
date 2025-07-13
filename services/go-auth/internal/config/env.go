package config

import (
	"github.com/caarlos0/env/v11"
)

type Envs struct {
	Secret           string `env:"JWT_SECRET"`
	RedisAddr        string `env:"REDIS_ADDRESS"`
	MetricsAddr      string `env:"METRICS_ADDRESS"`
	ServerAddr       string `env:"SERVER_ADDRESS"`
	JaegerAddr       string `env:"JEAGER_ADDRESS"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDB       string `env:"POSTGRES_DB"`
	MigrationPath    string `env:"MIGRATIONS_PATH"`
}

func ParseEnv() (*Envs, error) {
	e := Envs{}
	if err := env.Parse(&e); err != nil {
		return nil, err
	}
	return &e, nil
}
