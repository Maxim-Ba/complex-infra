package config

import "strings"

type Config struct {
	RedisAddr              string
	MetricsAddr            string
	ServerAddr             string
	JaegerAddr             string
	PostgresHost           string
	PostgresPort           string
	PostgresUser           string
	PostgresPassword       string
	PostgresDB             string
	MigrationPath          string
	KafkaAddr              string
	KafkaBrokers           []string
	KafkaGroupId           string
	RTCSignalTopic         string
	RTCResponseTopic       string
	WebRTCIceServers       []string // STUN/TURN серверы
	WebRTCSignalingTimeout int32    // Таймаут сигналинга в секундах
	ExternalIP             string
}

func New() *Config {
	cfg, err := ParseEnv()
	if err != nil {
		panic(err.Error())
	}
	return &Config{
		RedisAddr:              cfg.RedisAddr,
		MetricsAddr:            cfg.MetricsAddr,
		ServerAddr:             cfg.ServerAddr,
		JaegerAddr:             cfg.JaegerAddr,
		PostgresHost:           cfg.PostgresHost,
		PostgresPort:           cfg.PostgresPort,
		PostgresUser:           cfg.PostgresUser,
		PostgresPassword:       cfg.PostgresPassword,
		PostgresDB:             cfg.PostgresDB,
		MigrationPath:          cfg.MigrationPath,
		KafkaAddr:              cfg.KafkaAddr,
		KafkaBrokers:           strings.Split(cfg.KafkaBrokers, ","),
		KafkaGroupId:           cfg.KafkaGroupId,
		RTCSignalTopic:         cfg.RTCSignalTopic,
		RTCResponseTopic:       cfg.RTCResponseTopic,
		WebRTCIceServers:       cfg.WebRTCIceServers,
		WebRTCSignalingTimeout: cfg.WebRTCSignalingTimeout,
		ExternalIP:             cfg.ExternalIP,
	}
}
func (cfg *Config) GetConfig() *Config {
	return cfg
}
