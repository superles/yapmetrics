package config

var AgentConfig Config

func InitConfig() {

	flagConfig := parseFlags()
	envConfig := parseEnv()

	if len(envConfig.Endpoint) > 0 {
		AgentConfig.Endpoint = envConfig.Endpoint
	} else {
		AgentConfig.Endpoint = flagConfig.Endpoint
	}

	if envConfig.ReportInterval > 0 {
		AgentConfig.ReportInterval = envConfig.ReportInterval
	} else {
		AgentConfig.ReportInterval = flagConfig.ReportInterval
	}

	if envConfig.PollInterval > 0 {
		AgentConfig.PollInterval = envConfig.PollInterval
	} else {
		AgentConfig.PollInterval = flagConfig.PollInterval
	}
}
