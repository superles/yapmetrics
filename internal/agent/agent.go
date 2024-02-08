package agent

import (
	"context"
	"errors"
	"github.com/avast/retry-go/v4"
	"github.com/superles/yapmetrics/internal/agent/client"
	"github.com/superles/yapmetrics/internal/agent/config"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"net"
	"time"
)

const attempts = 4

type Agent struct {
	storage metricProvider
	config  *config.Config
	client  client.Client
}

// New Создание нового агента.
func New(s metricProvider, cfg *config.Config, cl client.Client) *Agent {
	agent := &Agent{storage: s, config: cfg, client: cl}
	return agent
}

func (a *Agent) sendAll(ctx context.Context, metrics []types.Metric) error {

	err := retry.Do(
		func() error {
			err := a.client.Send(ctx, a.config.Endpoint, metrics)
			if err != nil {
				var opError *net.OpError
				if errors.As(err, &opError) {
					return err
				} else {
					return retry.Unrecoverable(err)
				}
			}

			return nil
		},
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			delay := int(n*2 + 1)
			return time.Duration(delay) * time.Second
		}),
		retry.Attempts(uint(attempts)),
	)

	return err
}

func (a *Agent) poolTickRuntime(ctx context.Context, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logger.Log.Debug("captureRuntime")
			err := a.captureRuntime(ctx)
			if err != nil {
				logger.Log.Error("captureRuntime error", err.Error())
			}
		}
	}
}

func (a *Agent) poolTickPsutil(ctx context.Context, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logger.Log.Debug("capturePsutil")
			err := a.capturePsutil(ctx)
			if err != nil {
				logger.Log.Error("capturePsutil error", err.Error())
			}
		}
	}
}

func (a *Agent) sendTicker(ctx context.Context, reportInterval time.Duration) {
	ticker := time.NewTicker(reportInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logger.Log.Debug("agent run sendAllJSON")
			var metrics types.Collection
			var err error
			metrics, err = a.storage.GetAll(ctx)
			if err != nil {
				logger.Log.Error("storage get all error", err.Error())
				continue
			}
			err = a.sendAll(ctx, metrics.ToSlice())
			if err != nil {
				logger.Log.Error("sendAll error", err.Error())
			}
		}
	}
}

// Run Запуск агента.
func (a *Agent) Run(ctx context.Context) error {

	logger.Log.Debug("agent run")

	reportInterval := time.Second * time.Duration(a.config.ReportInterval)
	pollInterval := time.Second * time.Duration(a.config.PollInterval)

	go a.poolTickRuntime(ctx, pollInterval)
	go a.poolTickPsutil(ctx, pollInterval)

	if a.config.RateLimit > 0 {
		a.sendPoolTicker(ctx, reportInterval)
	} else {
		go a.sendTicker(ctx, reportInterval)
	}

	return nil
}
