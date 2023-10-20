package main

import (
	"github.com/superles/yapmetrics/internal/server"
	"github.com/superles/yapmetrics/internal/server/config"
	"github.com/superles/yapmetrics/internal/storage/pgstorage"
)

func main() {
	cfg := config.New()
	storage := pgstorage.New(cfg.DatabaseDsn)
	srv := server.New(storage, cfg)
	srv.Run()
}
