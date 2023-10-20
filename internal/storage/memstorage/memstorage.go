package memstorage

import (
	"context"
	types "github.com/superles/yapmetrics/internal/metric"
	"sync"
)

var storageSync = sync.Mutex{}

type MemStorage struct {
	collection map[string]types.Metric
}

func New() *MemStorage {
	return &MemStorage{make(map[string]types.Metric)}
}

func (s *MemStorage) GetAll(ctx context.Context) map[string]types.Metric {
	storageSync.Lock()
	defer storageSync.Unlock()
	targetMap := make(map[string]types.Metric, len(s.collection))

	// Copy from the original map to the target map
	for key, value := range s.collection {
		targetMap[key] = value
	}

	return targetMap
}

func (s *MemStorage) Get(ctx context.Context, name string) (types.Metric, bool) {
	storageSync.Lock()
	defer storageSync.Unlock()
	val, ok := s.collection[name]
	if !ok {
		return types.Metric{}, false
	}
	return val, true
}

func (s *MemStorage) Set(ctx context.Context, data *types.Metric) {
	storageSync.Lock()
	defer storageSync.Unlock()
	s.collection[data.Name] = *data
}

func (s *MemStorage) SetAll(ctx context.Context, data *[]types.Metric) error {
	storageSync.Lock()
	defer storageSync.Unlock()

	// Copy from the original map to the target map
	for _, value := range *data {
		switch value.Type {
		case types.GaugeMetricType:
			s.SetFloat(ctx, value.Name, value.Value)
		case types.CounterMetricType:
			s.IncCounter(ctx, value.Name, int64(value.Value))
		}
	}

	return nil
}

func (s *MemStorage) SetFloat(ctx context.Context, Name string, Value float64) {
	storageSync.Lock()
	defer storageSync.Unlock()
	s.collection[Name] = types.Metric{Name: Name, Type: types.GaugeMetricType, Value: Value}
}

func (s *MemStorage) IncCounter(ctx context.Context, Name string, Value int64) {
	storageSync.Lock()
	defer storageSync.Unlock()
	val := s.collection[Name]
	val.Name = Name
	val.Type = types.CounterMetricType
	val.Value = val.Value + float64(Value)
	s.collection[Name] = val
}
