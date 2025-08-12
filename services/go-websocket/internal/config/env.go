package config

import (
	"github.com/caarlos0/env/v11"
)

type Envs struct {
	MetricsAddr               string `env:"METRICS_ADDRESS"`
	ServerAddr                string `env:"SERVER_ADDRESS"`
	JaegerAddr                string `env:"JEAGER_ADDRESS"`
	MessageTopic              string `env:"MESSAGE_TOPIC"`
	KafkaAddr                 string `env:"KAFKA_ADDRESS"`
	KafkaBrokers              string `env:"KAFKA_BROKERS"`
	KafkaGroupId              string `env:"KAFKA_GROUP_ID"`
	MessageConfirmationsTopic string `env:"MESSAGE_CONFIRMATIONS_TOPIC"`
}

func ParseEnv() (*Envs, error) {
	e := Envs{}
	if err := env.Parse(&e); err != nil {
		return nil, err
	}
	return &e, nil
}
