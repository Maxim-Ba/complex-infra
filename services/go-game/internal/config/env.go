package config

import (
	"github.com/caarlos0/env/v11"
)

type Envs struct {
	RedisAddr              string   `env:"REDIS_ADDRESS"`
	MetricsAddr            string   `env:"METRICS_ADDRESS"`
	ServerAddr             string   `env:"SERVER_ADDRESS"`
	JaegerAddr             string   `env:"JEAGER_ADDRESS"`
	PostgresHost           string   `env:"POSTGRES_HOST"`
	PostgresPort           string   `env:"POSTGRES_PORT"`
	PostgresUser           string   `env:"POSTGRES_USER"`
	PostgresPassword       string   `env:"POSTGRES_PASSWORD"`
	PostgresDB             string   `env:"POSTGRES_DB"`
	MigrationPath          string   `env:"MIGRATIONS_PATH"`
	KafkaAddr              string   `env:"KAFKA_ADDRESS"`
	KafkaBrokers           string   `env:"KAFKA_BROKERS"`
	KafkaGroupId           string   `env:"KAFKA_GROUP_ID"`
	RTCSignalTopic         string   `env:"RTC_SIGNAL_TOPIC"`
	RTCResponseTopic       string   `env:"RTC_RESPONSE_TOPIC"`
	WebRTCIceServers       []string `env:"RTC_ICE_SERVERS" envSeparator:","`
	WebRTCSignalingTimeout int32    `env:"RTC_SIGNAL_TIMEOUT"`
	ExternalIP             string   `env:"EXTERNAL_IP"`
}

func ParseEnv() (*Envs, error) {
	e := Envs{}
	if err := env.Parse(&e); err != nil {
		return nil, err
	}
	return &e, nil
}
