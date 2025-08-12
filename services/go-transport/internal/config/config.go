package config

type Config struct {
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
		MetricsAddr: cfg.MetricsAddr,
		ServerAddr:  cfg.ServerAddr,
		JaegerAddr:  cfg.JaegerAddr,
	}
}
func (cfg *Config) GetConfig() *Config {
	return cfg
}
