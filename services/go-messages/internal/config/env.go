package config

import (
	"github.com/caarlos0/env/v11"
)

type Envs struct {
	Secret          string `env:"JWT_SECRET"`
	RedisAddr       string `env:"REDIS_ADDRESS"`
	MetricsAddr     string `env:"METRICS_ADDRESS"`
	ServerAddr      string `env:"SERVER_ADDRESS"`
	JaegerAddr      string `env:"JEAGER_ADDRESS"`
	KafkaAddr       string `env:"KAFKA_ADDRESS"`
	KafkaBrokers    string `env:"KAFKA_BROKERS"`
	KafkaGroupId    string `env:"KAFKA_GROUP_ID"`
	MongoDBURI      string `env:"MONGODB_URI"`
	MongoDBDatabase string `env:"MONGODB_DATABASE"`
}

func ParseEnv() (*Envs, error) {
	e := Envs{}
	if err := env.Parse(&e); err != nil {
		return nil, err
	}
	return &e, nil
}
