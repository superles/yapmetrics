package main

import (
	"github.com/superles/yapmetrics/internal/server"
	"github.com/superles/yapmetrics/internal/server/config"
	storage2 "github.com/superles/yapmetrics/internal/storage"
	"github.com/superles/yapmetrics/internal/storage/memstorage"
	"github.com/superles/yapmetrics/internal/storage/pgstorage"
)

func main() {
	cfg := config.New()
	var storage storage2.Storage
	storage = memstorage.New()
	if len(cfg.DatabaseDsn) > 0 {
		storage = pgstorage.New(cfg.DatabaseDsn)
	}
	srv := server.New(storage, cfg)
	srv.Run()
}
