package repository

import (
	"errors"
	"github.com/superles/yapmetrics/internal/types"
)

type MemoryMetricRepository struct {
	store types.MetricCollection
}

func (r *MemoryMetricRepository) Set(Name string, Type string, Value interface{}) {
	if r.store == nil {
		r.store = types.MetricCollection{}
	}
	r.store[Name] = types.Metric{
		Name:  Name,
		Type:  Type,
		Value: Value,
	}
}

func (r *MemoryMetricRepository) GetAll() types.MetricCollection {
	if r.store == nil || len(r.store) == 0 {
		return types.MetricCollection{}
	}
	targetMap := make(types.MetricCollection)

	// Copy from the original map to the target map
	for key, value := range r.store {
		targetMap[key] = value
	}

	return targetMap
}

func (r *MemoryMetricRepository) Get(name string) (types.Metric, error) {
	if r.store == nil || len(r.store) == 0 {
		return types.Metric{}, errors.New("метрика не найдена")
	}
	if _, ok := r.store[name]; !ok {
		return types.Metric{}, errors.New("метрика не найдена")
	}
	return r.store[name], nil
}
