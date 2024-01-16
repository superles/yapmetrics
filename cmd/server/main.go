package main

// @title           Swagger Server API
// @version         1.0
// @description     Yandex Practicum metrics server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

import (
	"context"
	"github.com/superles/yapmetrics/internal/server"
	"github.com/superles/yapmetrics/internal/server/config"
	"github.com/superles/yapmetrics/internal/storage"
	"github.com/superles/yapmetrics/internal/storage/memstorage"
	"github.com/superles/yapmetrics/internal/storage/pgstorage"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"log"
)

func main() {
	cfg := config.New()
	var store storage.Storage
	var err error
	if len(cfg.DatabaseDsn) != 0 {
		if store, err = pgstorage.New(cfg.DatabaseDsn); err != nil {
			log.Fatal("ошибка инициализации бд", err.Error())
		}
	} else {
		store = memstorage.New()
	}
	if err = logger.Initialize(cfg.LogLevel); err != nil {
		log.Fatal("ошибка инициализации логера", err.Error())
	}
	srv := server.New(store, cfg)
	appContext := context.Background()
	if err = srv.Run(appContext); err != nil {
		log.Fatal("ошибка запуска сервера", err.Error())
	}
}
