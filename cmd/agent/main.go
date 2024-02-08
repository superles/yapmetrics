package main

import (
	"context"
	"fmt"
	"github.com/superles/yapmetrics/internal/agent/client"
	grpc "github.com/superles/yapmetrics/internal/grpc/client"
	"github.com/superles/yapmetrics/internal/utils/encoder"
	"github.com/superles/yapmetrics/internal/utils/network"
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

	systemIP := cfg.RealIP

	if len(systemIP) == 0 {
		systemIP, err = network.ParseIP()
		if err != nil {
			// error of catching ip is not critical
			logger.Log.Error(err.Error())
		}
	}

	var enc *encoder.Encoder

	if len(cfg.CryptoKey) != 0 {
		enc, err = encoder.New(cfg.CryptoKey)
		if err != nil {
			log.Panicln("ошибка инициализации encoder", err.Error())
		}
	}

	var cl client.Client

	if cfg.ClientType == "grpc" {
		cl = grpc.NewGrpcClient(grpc.GrpcClientParams{Key: cfg.SecretKey, RealIP: systemIP, Encoder: enc})
	} else {
		params := client.AgentClientParams{Key: cfg.SecretKey, RealIP: systemIP, Compress: true, Encoder: enc}
		cl = client.NewHTTPAgentClient(params)
	}

	app := agent.New(storage, cfg, cl, enc)

	if err = app.Run(appContext); err != nil {
		log.Panic("ошибка запуска агента", err.Error())
	}
	logger.Log.Info("Agent Started")
	<-appContext.Done()
	logger.Log.Info("Agent Stopped")
	logger.Log.Info("Agent Exited Properly")
}
