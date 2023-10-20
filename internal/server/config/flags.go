package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func parseFlags() Config {

	var config Config

	flag.StringVar(&config.Endpoint, "a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	flag.StringVar(&config.LogLevel, "v", "debug", "уровень логирования")
	flag.IntVar(&config.StoreInterval, "i", 300, "интервал сохранения метрик на диск")
	flag.StringVar(&config.FileStoragePath, "f", filepath.Join(os.TempDir(), "metrics-db.json"), "интервал сохранения метрик на диск")
	flag.BoolVar(&config.Restore, "r", true, "интервал сохранения метрик на диск")
	//example: postgresql://test_user:test_user@localhost/test_db
	flag.StringVar(&config.DatabaseDsn, "d", "", "интервал сохранения метрик на диск")

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
