package config

import (
	"flag"
	"fmt"
)

func parseFlags() Config {

	var config Config

	flag.StringVar(&config.Endpoint, "a", "localhost:8080", "адрес эндпоинта HTTP-сервера")

	var Usage = func() {
		_, err := fmt.Fprintf(flag.CommandLine.Output(), "Параметры командной строки сервера:\n")
		if err != nil {
			fmt.Println(err.Error())
		}
		flag.PrintDefaults()
	}
	flag.Usage = Usage
	flag.Parse()

	return config
}
