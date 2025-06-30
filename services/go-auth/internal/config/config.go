package config

type Config struct {
	Secret      string
	RedisAddr   string
	MetricsAddr string
	ServerAddr  string
	JaegerAddr  string
}

func New() *Config {
	cfg, err := ParseEnv()
	if err != nil {
		panic(err.Error())
	}
	return &Config{
		Secret:      cfg.Secret,
		RedisAddr:   cfg.RedisAddr,
		MetricsAddr: cfg.MetricsAddr,
		ServerAddr:  cfg.ServerAddr,
		JaegerAddr:  cfg.JaegerAddr,
	}
}
