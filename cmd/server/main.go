package main

import (
	"github.com/superles/yapmetrics/internal/server"
	"github.com/superles/yapmetrics/internal/server/config"
	"github.com/superles/yapmetrics/internal/storage"
	"github.com/superles/yapmetrics/internal/storage/memstorage"
	"github.com/superles/yapmetrics/internal/storage/pgstorage"
	"github.com/superles/yapmetrics/internal/utils/logger"
)

func main() {
	cfg := config.New()
	var store storage.Storage
	store = memstorage.New()
	if len(cfg.DatabaseDsn) > 0 {
		logger.Log.Debug(cfg)
		store = pgstorage.New(cfg.DatabaseDsn)
	}
	srv := server.New(store, cfg)
	srv.Run()
}
