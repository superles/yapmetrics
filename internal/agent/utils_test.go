package agent

import (
	"github.com/stretchr/testify/require"
	"github.com/superles/yapmetrics/internal/metric"
	"testing"
)

func Test_compressMetrics(t *testing.T) {
	metrics := make(map[string]metric.Metric, 1)
	t.Run("test empty compressMetrics", func(t *testing.T) {
		_, err := compressMetrics(metrics)
		require.NoError(t, err)
	})
	metrics["test"] = metric.Metric{Name: "test", Type: metric.CounterMetricType, Value: 1}
	t.Run("test full compressMetrics", func(t *testing.T) {
		_, err := compressMetrics(metrics)
		require.NoError(t, err)
	})
}
