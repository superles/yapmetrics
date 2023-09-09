package repository

import "github.com/superles/yapmetrics/internal/types"

type MetricRepositoryInterface interface {
	GetAll() types.MetricCollection
	Get(name string) (types.Metric, error)
	Set(Name string, Type string, Value interface{})
}
