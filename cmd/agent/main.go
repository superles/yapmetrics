package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/superles/yapmetrics/internal/agent"
	"github.com/superles/yapmetrics/internal/agent/config"
	"github.com/superles/yapmetrics/internal/storage/memstorage"
	"github.com/superles/yapmetrics/internal/utils/logger"
)

func main() {
	storage := memstorage.New()
	cfg := config.New()
	err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Panicln("ошибка инициализации логера", err.Error())
	}

	appContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err = agent.New(storage, cfg).Run(appContext); err != nil {
		log.Fatal("ошибка запуска агента", err.Error())
	}
	logger.Log.Info("Agent Started")
	<-appContext.Done()
	logger.Log.Info("Agent Stopped")
	logger.Log.Info("Agent Exited Properly")
}
