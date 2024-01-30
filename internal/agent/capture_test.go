package agent

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/superles/yapmetrics/internal/agent/client"
	"github.com/superles/yapmetrics/internal/agent/config"
	"github.com/superles/yapmetrics/internal/storage/memstorage"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"testing"
)

func Benchmark_captureRuntime(b *testing.B) {

	cfg := config.New()
	agent := New(memstorage.New(), config.New(), client.NewHTTPAgentClient(client.AgentClientParams{Key: cfg.SecretKey}), nil)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := agent.captureRuntime(ctx)
		if err != nil {
			b.Fatalf("Failed to capture runtime: %v", err)
		}
	}
}

func Benchmark_capturePsutil(b *testing.B) {

	cfg := config.New()
	agent := New(memstorage.New(), config.New(), client.NewHTTPAgentClient(client.AgentClientParams{Key: cfg.SecretKey}), nil)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := agent.capturePsutil(ctx)
		if err != nil {
			b.Fatalf("Failed to capture psutil metrics: %v", err)
		}
	}
}

func TestAgent_captureRuntime(t *testing.T) {
	t.Run("test captureRuntime", func(t *testing.T) {
		cfg := config.New()
		storage := memstorage.New()
		err := logger.Initialize(cfg.LogLevel)
		require.NoError(t, err)

		a := New(storage, config.New(), client.NewHTTPAgentClient(client.AgentClientParams{Key: cfg.SecretKey}), nil)

		err = a.captureRuntime(context.Background())
		require.NoError(t, err)
	})
}

func TestAgent_capturePsutil(t *testing.T) {
	t.Run("test capturePsutil", func(t *testing.T) {
		cfg := config.New()
		storage := memstorage.New()
		err := logger.Initialize(cfg.LogLevel)
		require.NoError(t, err)

		a := New(storage, config.New(), client.NewHTTPAgentClient(client.AgentClientParams{Key: cfg.SecretKey}), nil)

		err = a.capturePsutil(context.Background())
		require.NoError(t, err)
	})
}

func Test_randFloat(t *testing.T) {
	t.Run("test randFloat", func(t *testing.T) {
		if got := randFloat(0, 100); got <= 0 || got >= 100 {
			t.Errorf("randFloat() = %v, want 0<%v<100", got, got)
		}
	})
}
