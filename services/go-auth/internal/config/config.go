package config

type Config struct {
	Secret    string
	RedisAddr string
}

func New() *Config {
	cfg, err := ParseEnv()
	if err != nil {
		panic(err.Error())
	}
	return &Config{
		Secret:    cfg.Secret,
		RedisAddr: cfg.RedisAddr,
	}
}
