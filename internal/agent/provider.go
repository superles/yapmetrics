package agent

import (
	"context"
	types "github.com/superles/yapmetrics/internal/metric"
)

type metricProvider interface {
	GetAll(ctx context.Context) (map[string]types.Metric, error)
	SetFloat(ctx context.Context, Name string, Value float64) error
	IncCounter(ctx context.Context, Name string, Value int64) error
}
