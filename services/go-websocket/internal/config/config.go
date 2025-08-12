package config

import "strings"

type Config struct {
	MetricsAddr               string
	ServerAddr                string
	JaegerAddr                string
	MessageTopic              string
	KafkaAddr                 string
	KafkaBrokers              []string
	KafkaGroupId              string
	MessageConfirmationsTopic string
}

func New() *Config {
	cfg, err := ParseEnv()
	if err != nil {
		panic(err.Error())
	}
	return &Config{
		MetricsAddr:               cfg.MetricsAddr,
		ServerAddr:                cfg.ServerAddr,
		JaegerAddr:                cfg.JaegerAddr,
		MessageTopic:              cfg.MessageTopic,
		KafkaAddr:                 cfg.KafkaAddr,
		KafkaBrokers:              strings.Split(cfg.KafkaBrokers, ","),
		KafkaGroupId:              cfg.KafkaGroupId,
		MessageConfirmationsTopic: cfg.MessageConfirmationsTopic,
	}
}
func (cfg *Config) GetConfig() *Config {
	return cfg
}
