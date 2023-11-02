package main

import (
	"context"
	"github.com/superles/yapmetrics/internal/agent"
	"github.com/superles/yapmetrics/internal/agent/config"
	"github.com/superles/yapmetrics/internal/storage/memstorage"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"log"
)

func main() {
	storage := memstorage.New()
	cfg := config.New()
	err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Panicln("ошибка инициализации логера", err.Error())
	}

	appContext := context.Background()
	if err = agent.New(storage, cfg).Run(appContext); err != nil {
		log.Fatal("ошибка запуска агента", err.Error())
	}

}
