package config

import (
	"flag"
	"fmt"
)

func parseFlags() Config {

	var config Config

	flag.StringVar(&config.Endpoint, "a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	flag.IntVar(&config.ReportInterval, "r", 10, "частота отправки метрик на сервер")
	flag.IntVar(&config.PollInterval, "p", 2, "частота опроса метрик из пакета runtime")
	flag.StringVar(&config.LogLevel, "v", "info", "уровень логирования")
	flag.StringVar(&config.SecretKey, "k", "", "Секретный ключ для хеширования ответов и проверки запросов")
	flag.UintVar(&config.RateLimit, "l", 0, "Количество одновременно исходящих запросов на сервер")

	var Usage = func() {
		_, err := fmt.Fprintf(flag.CommandLine.Output(), "Параметры командной строки агента:\n")
		if err != nil {
			fmt.Println(err.Error())
		}
		flag.PrintDefaults()
	}
	flag.Usage = Usage
	flag.Parse()

	return config
}
