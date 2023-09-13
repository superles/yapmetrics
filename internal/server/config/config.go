package config

func Load() *Config {

	var ServerConfig Config

	flagConfig := parseFlags()
	envConfig := parseEnv()

	if len(envConfig.Endpoint) > 0 {
		ServerConfig.Endpoint = envConfig.Endpoint
	} else {
		ServerConfig.Endpoint = flagConfig.Endpoint
	}

	return &ServerConfig
}
