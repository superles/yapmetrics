package agent

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/superles/yapmetrics/internal/agent/client"
	"github.com/superles/yapmetrics/internal/agent/config"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/storage/memstorage"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"testing"
	"time"
)

func TestAgent_generator(t *testing.T) {
	store := memstorage.New()
	cfg := config.New()
	agentClient := client.NewHTTPAgentClient(client.AgentClientParams{Key: cfg.SecretKey})
	err := logger.Initialize(cfg.LogLevel)
	require.NoError(t, err)
	a := Agent{storage: store, config: cfg, client: agentClient}
	ctx, done := context.WithCancel(context.Background())
	defer done()
	err = store.Set(ctx, types.Metric{Name: "test", Type: types.CounterMetricType, Value: 1})
	require.NoError(t, err)
	requestChan := make(chan types.Collection, 1)
	defer close(requestChan)
	t.Run("generator test", func(t *testing.T) {
		go a.generator(context.Background(), requestChan, 100*time.Millisecond)
		select {
		case <-ctx.Done():
			return // Выход из горутины при отмене контекста
		case _, ok := <-requestChan:
			require.Exactly(t, ok, true, "generator not generate")
		}
	})
}
