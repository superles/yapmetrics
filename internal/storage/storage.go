package storage

import (
	"context"
	types "github.com/superles/yapmetrics/internal/metric"
)

type Counter int64
type Gauge float64

type Storage interface {
	GetAll(ctx context.Context) map[string]types.Metric
	Get(ctx context.Context, name string) (types.Metric, bool)
	Set(ctx context.Context, data *types.Metric)
	SetFloat(ctx context.Context, Name string, Value float64)
	IncCounter(ctx context.Context, Name string, Value int64)
}
