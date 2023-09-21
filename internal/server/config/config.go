package config

import "sync"

type Config struct {
	Endpoint string `env:"ADDRESS"`
	LogLevel string `env:"SERVER_LOG_LEVEL"`
}

var (
	once     sync.Once
	instance Config
)

func New() *Config {

	once.Do(func() {

		instance = Config{}

		flagConfig := parseFlags()
		envConfig := parseEnv()

		if len(envConfig.Endpoint) > 0 {
			instance.Endpoint = envConfig.Endpoint
		} else {
			instance.Endpoint = flagConfig.Endpoint
		}

		if len(envConfig.LogLevel) > 0 {
			instance.LogLevel = envConfig.LogLevel
		} else {
			instance.LogLevel = flagConfig.LogLevel
		}
	})

	return &instance
}
