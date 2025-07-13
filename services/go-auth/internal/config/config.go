package config

type Config struct {
	Secret           string
	RedisAddr        string
	MetricsAddr      string
	ServerAddr       string
	JaegerAddr       string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	MigrationPath    string
}

func New() *Config {
	cfg, err := ParseEnv()
	if err != nil {
		panic(err.Error())
	}
	return &Config{
		Secret:           cfg.Secret,
		RedisAddr:        cfg.RedisAddr,
		MetricsAddr:      cfg.MetricsAddr,
		ServerAddr:       cfg.ServerAddr,
		JaegerAddr:       cfg.JaegerAddr,
		PostgresHost:     cfg.PostgresHost,
		PostgresPort:     cfg.PostgresPort,
		PostgresUser:     cfg.PostgresUser,
		PostgresPassword: cfg.PostgresPassword,
		PostgresDB:       cfg.PostgresDB,
		MigrationPath:    cfg.MigrationPath,
	}
}
func (cfg *Config) GetConfig() *Config {
	return cfg
}
