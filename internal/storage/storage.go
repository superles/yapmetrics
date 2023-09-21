package storage

import (
	types "github.com/superles/yapmetrics/internal/metric"
)

type Counter int64
type Gauge float64

type Storage interface {
	GetAll() map[string]types.Metric
	Get(name string) (types.Metric, bool)
	Set(data *types.Metric)
	SetFloat(Name string, Value float64)
	IncCounter(Name string, Value int64)
}
