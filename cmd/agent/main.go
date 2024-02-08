package main

import (
	"context"
	"fmt"
	"github.com/superles/yapmetrics/internal/agent/client"
	"github.com/superles/yapmetrics/internal/utils/encoder"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/superles/yapmetrics/internal/agent"
	"github.com/superles/yapmetrics/internal/agent/config"
	"github.com/superles/yapmetrics/internal/storage/memstorage"
	"github.com/superles/yapmetrics/internal/utils/logger"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func printInfo() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func main() {
	printInfo()
	storage := memstorage.New()
	cfg := config.New()
	err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Panicln("ошибка инициализации логера", err.Error())
	}

	appContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cl := client.NewHTTPAgentClient(client.AgentClientParams{Key: cfg.SecretKey})

	var app *agent.Agent

	if len(cfg.CryptoKey) != 0 {
		enc, err := encoder.New(cfg.CryptoKey)
		if err != nil {
			log.Panicln("ошибка инициализации encoder", err.Error())
		}
		app = agent.New(storage, cfg, cl, enc)
	} else {
		app = agent.New(storage, cfg, cl, nil)
	}

	if err = app.Run(appContext); err != nil {
		log.Fatal("ошибка запуска агента", err.Error())
	}
	logger.Log.Info("Agent Started")
	<-appContext.Done()
	logger.Log.Info("Agent Stopped")
	logger.Log.Info("Agent Exited Properly")
}
