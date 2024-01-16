package config

import "sync"

type Config struct {
	Endpoint        string `env:"ADDRESS"`
	LogLevel        string `env:"SERVER_LOG_LEVEL"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDsn     string `env:"DATABASE_DSN"`
	SecretKey       string `env:"KEY"`
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

		if envConfig.StoreInterval > 0 {
			instance.StoreInterval = envConfig.StoreInterval
		} else {
			instance.StoreInterval = flagConfig.StoreInterval
		}

		if len(envConfig.FileStoragePath) > 0 {
			instance.FileStoragePath = envConfig.FileStoragePath
		} else {
			instance.FileStoragePath = flagConfig.FileStoragePath
		}

		if envConfig.Restore {
			instance.Restore = envConfig.Restore
		} else {
			instance.Restore = flagConfig.Restore
		}

		if len(envConfig.DatabaseDsn) > 0 {
			instance.DatabaseDsn = envConfig.DatabaseDsn
		} else {
			instance.DatabaseDsn = flagConfig.DatabaseDsn
		}

		if len(envConfig.SecretKey) > 0 {
			instance.SecretKey = envConfig.SecretKey
		} else {
			instance.SecretKey = flagConfig.SecretKey
		}
	})

	return &instance
}
