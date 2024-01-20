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
	agentClient := client.NewHTTPAgentClient(client.AgentClientParams{})
	err := logger.Initialize(cfg.LogLevel)
	require.NoError(t, err)
	a := Agent{store, cfg, agentClient, logger.Log}
	ctx, done := context.WithCancel(context.Background())
	defer done()
	err = store.Set(ctx, types.Metric{Name: "test", Type: types.CounterMetricType, Value: 1})
	require.NoError(t, err)
	input := make(chan<- types.Collection, 3)
	defer close(input)
	t.Run("generator test", func(t *testing.T) {
		go a.generator(context.Background(), input, 100*time.Millisecond)
		time.Sleep(1000 * time.Millisecond)
		require.Exactly(t, len(input), 3, "wrong generator elements count")
	})
}
