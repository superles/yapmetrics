package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Config struct {
	Endpoint        string `env:"ADDRESS" json:"address"`
	LogLevel        string `env:"SERVER_LOG_LEVEL"`
	StoreInterval   int    `env:"STORE_INTERVAL" json:"store_interval"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"store_file"`
	Restore         bool   `env:"RESTORE" json:"restore"`
	DatabaseDsn     string `env:"DATABASE_DSN" json:"database_dsn"`
	SecretKey       string `env:"KEY"`
	CryptoKey       string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigFile      string `env:"CONFIG"`
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

		if len(envConfig.ConfigFile) > 0 {
			instance.ConfigFile = envConfig.ConfigFile
		} else {
			instance.ConfigFile = flagConfig.ConfigFile
		}

		if len(instance.ConfigFile) > 0 {
			instance = parseFile(instance.ConfigFile)
			return
		}

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

		if len(envConfig.CryptoKey) > 0 {
			instance.CryptoKey = envConfig.CryptoKey
		} else {
			instance.CryptoKey = flagConfig.CryptoKey
		}
	})

	return &instance
}

func parseFile(fileName string) Config {

	data, err := os.ReadFile(fileName)

	if err != nil {
		log.Fatal(err)
	}

	var config Config

	err = json.Unmarshal(data, &config)

	if err != nil {
		log.Fatal(err)
	}

	return config
}
