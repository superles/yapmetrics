package storage

import (
	"context"
	types "github.com/superles/yapmetrics/internal/metric"
)

type Counter int64
type Gauge float64

type Storage interface {
	GetAll(ctx context.Context) (map[string]types.Metric, error)
	Get(ctx context.Context, name string) (types.Metric, error)
	Set(ctx context.Context, data *types.Metric) error
	SetAll(ctx context.Context, data *[]types.Metric) error
	SetFloat(ctx context.Context, Name string, Value float64) error
	IncCounter(ctx context.Context, Name string, Value int64) error
	Ping(ctx context.Context) error
}
