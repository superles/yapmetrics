package storage

import (
	types "github.com/superles/yapmetrics/internal/metric"
)

type Counter int64
type Gauge float64

type Storage interface {
	GetAll() map[string]types.Metric
	Get(name string) (types.Metric, error)
	SetFloat(Name string, Value float64)
	SetInt(Name string, Value int64)
	IncCounter(Name string, Value int64)
}
