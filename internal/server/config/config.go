package config

var ServerConfig Config

func InitConfig() {

	flagConfig := parseFlags()
	envConfig := parseEnv()

	if len(envConfig.Endpoint) > 0 {
		ServerConfig.Endpoint = envConfig.Endpoint
	} else {
		ServerConfig.Endpoint = flagConfig.Endpoint
	}
}
