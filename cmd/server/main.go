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
	if len(cfg.DatabaseDsn) != 0 {
		store = pgstorage.New(cfg.DatabaseDsn)
	} else {
		store = memstorage.New()
	}
	err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Panicln("ошибка инициализации логера", err.Error())
	}
	srv := server.New(store, cfg)
	srv.Run()
}
