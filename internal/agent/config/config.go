package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type Config struct {
	Endpoint       string `env:"ADDRESS" json:"address"`
	ReportInterval int    `env:"REPORT_INTERVAL" json:"report_interval"`
	PollInterval   int    `env:"POLL_INTERVAL" json:"poll_interval"`
	LogLevel       string `env:"AGENT_LOG_LEVEL"`
	SecretKey      string `env:"KEY"`
	RateLimit      uint   `env:"RATE_LIMIT"`
	CryptoKey      string `env:"CRYPTO_KEY" json:"crypto_key"`
	ConfigFile     string `env:"CONFIG"`
	RealIP         string `env:"REAL_IP" json:"real_ip"`
	ClientType     string `env:"CLIENT_TYPE" json:"client_type"` // exist types: http, grpc
}

var (
	once     sync.Once
	instance Config
)

func initConfig() {

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

	if envConfig.RateLimit > 0 {
		instance.RateLimit = envConfig.RateLimit
	} else {
		instance.RateLimit = flagConfig.RateLimit
	}

	if len(envConfig.RealIP) > 0 {
		instance.RealIP = envConfig.RealIP
	} else {
		instance.RealIP = flagConfig.RealIP
	}

	if len(envConfig.ClientType) > 0 {
		instance.ClientType = envConfig.ClientType
	} else {
		instance.ClientType = flagConfig.ClientType
	}

}

// New Создание объекта Config, pattern: Singleton.
func New() *Config {

	once.Do(initConfig)

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
