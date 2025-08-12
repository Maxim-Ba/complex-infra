package config

import (
	"github.com/caarlos0/env/v11"
)

type Envs struct {

	MetricsAddr      string `env:"METRICS_ADDRESS"`
	ServerAddr       string `env:"SERVER_ADDRESS"`
	JaegerAddr       string `env:"JEAGER_ADDRESS"`

}

func ParseEnv() (*Envs, error) {
	e := Envs{}
	if err := env.Parse(&e); err != nil {
		return nil, err
	}
	return &e, nil
}
