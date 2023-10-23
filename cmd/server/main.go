package main

import (
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
	srv.Run()
}
