package config

import (
	"strings"
)

type Config struct {
	Secret       string
	RedisAddr    string
	MetricsAddr  string
	ServerAddr   string
	JaegerAddr   string
	KafkaAddr    string
	KafkaBrokers []string
	KafkaGroupId string
}

func New() *Config {
	cfg, err := ParseEnv()
	if err != nil {
		panic(err.Error())
	}
	return &Config{
		Secret:       cfg.Secret,
		RedisAddr:    cfg.RedisAddr,
		MetricsAddr:  cfg.MetricsAddr,
		ServerAddr:   cfg.ServerAddr,
		JaegerAddr:   cfg.JaegerAddr,
		KafkaAddr:    cfg.KafkaAddr,
		KafkaBrokers: strings.Split(cfg.KafkaBrokers, ","),
		KafkaGroupId: cfg.KafkaGroupId,
	}
}
func (cfg *Config) GetConfig() *Config {
	return cfg
}
