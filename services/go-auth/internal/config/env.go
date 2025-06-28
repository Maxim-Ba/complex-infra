package config

import (
	"github.com/caarlos0/env/v11"
)

type Envs struct {
	Secret    string `env:"JWT_SECRET"`
	RedisAddr string `env:"REDIS_ADDRESS"`
}

func ParseEnv() (*Envs, error) {
	e := Envs{}
	if err := env.Parse(&e); err != nil {
		return nil, err
	}
	return &e, nil
}
