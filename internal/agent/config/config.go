package config

import "sync"

type Config struct {
	Endpoint       string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
	LogLevel       string `env:"AGENT_LOG_LEVEL"`
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

		if envConfig.ReportInterval > 0 {
			instance.ReportInterval = envConfig.ReportInterval
		} else {
			instance.ReportInterval = flagConfig.ReportInterval
		}

		if envConfig.PollInterval > 0 {
			instance.PollInterval = envConfig.PollInterval
		} else {
			instance.PollInterval = flagConfig.PollInterval
		}

		if len(envConfig.LogLevel) > 0 {
			instance.LogLevel = envConfig.LogLevel
		} else {
			instance.LogLevel = flagConfig.LogLevel
		}
	})

	return &instance
}
