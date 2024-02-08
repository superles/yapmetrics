package config

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
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

func (c *Config) merge(primaryConfig Config, secondaryConfig Config) {
	cfgType := reflect.TypeOf(c).Elem()
	cfgVal := reflect.ValueOf(c).Elem()
	primaryCfgVal := reflect.ValueOf(primaryConfig)
	secondaryCfgVal := reflect.ValueOf(secondaryConfig)
	fieldCount := cfgVal.NumField()
	for i := 0; i < fieldCount; i++ {
		field := cfgVal.Field(i)
		fieldName := cfgType.Field(i).Name
		switch field.Kind() {
		case reflect.String:
			val := field.String()
			if len(val) > 0 {
				continue
			}
			if primaryVal := primaryCfgVal.FieldByName(fieldName).String(); len(primaryVal) > 0 {
				field.SetString(primaryVal)
			} else if secondaryVal := secondaryCfgVal.FieldByName(fieldName).String(); len(secondaryVal) > 0 {
				field.SetString(secondaryVal)
			}
		case reflect.Int:
			val := field.Int()
			if val > 0 {
				continue
			}
			if primaryVal := primaryCfgVal.FieldByName(fieldName).Int(); primaryVal != 0 {
				field.SetInt(primaryVal)
			} else if secondaryVal := secondaryCfgVal.FieldByName(fieldName).Int(); secondaryVal != 0 {
				field.SetInt(secondaryVal)
			}
		case reflect.Uint:
			val := field.Uint()
			if val > 0 {
				continue
			}
			if primaryVal := primaryCfgVal.FieldByName(fieldName).Uint(); primaryVal > 0 {
				field.SetUint(primaryVal)
			} else if secondaryVal := secondaryCfgVal.FieldByName(fieldName).Uint(); secondaryVal > 0 {
				field.SetUint(secondaryVal)
			}
		case reflect.Bool:
			if primaryVal := primaryCfgVal.FieldByName(fieldName).Bool(); primaryVal {
				field.SetBool(primaryVal)
			} else if secondaryVal := secondaryCfgVal.FieldByName(fieldName).Bool(); secondaryVal {
				field.SetBool(secondaryVal)
			}
		default:
			panic("unhandled default case")
		}
	}
}

// New Создание объекта Config, pattern: Singleton.
func New() *Config {

	once.Do(func() {

		instance = Config{}

		flagConfig := parseFlags()
		envConfig := parseEnv()

		instance.merge(envConfig, flagConfig)

		if len(instance.ConfigFile) > 0 {
			instance = parseFile(instance.ConfigFile)
			return
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
