package main

import (
	"github.com/superles/yapmetrics/internal/server"
	"github.com/superles/yapmetrics/internal/storage/memstorage"
)

func main() {
	storage := memstorage.New()
	srv := server.New(storage)
	srv.Run()
}
