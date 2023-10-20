package main

import (
	"github.com/superles/yapmetrics/internal/agent"
	"github.com/superles/yapmetrics/internal/agent/config"
	"github.com/superles/yapmetrics/internal/storage/memstorage"
)

func main() {
	storage := memstorage.New()
	cfg := config.New()
	agent.New(storage, cfg).Run()
}
