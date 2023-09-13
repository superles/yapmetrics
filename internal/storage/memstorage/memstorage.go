package memstorage

import (
	"errors"
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

func (s *MemStorage) GetAll() map[string]types.Metric {
	storageSync.Lock()
	defer storageSync.Unlock()
	targetMap := make(map[string]types.Metric, len(s.collection))

	// Copy from the original map to the target map
	for key, value := range s.collection {
		targetMap[key] = value
	}

	return targetMap
}

func (s *MemStorage) Get(name string) (types.Metric, error) {
	storageSync.Lock()
	defer storageSync.Unlock()
	val, ok := s.collection[name]
	if !ok {
		return types.Metric{}, errors.New("метрика не найдена")
	}
	return val, nil
}

func (s *MemStorage) SetFloat(Name string, Value float64) {
	storageSync.Lock()
	defer storageSync.Unlock()
	s.collection[Name] = types.Metric{Name: Name, Type: types.GaugeMetricTypeName, ValueFloat: Value}
}

func (s *MemStorage) IncCounter(Name string, Value int64) {
	storageSync.Lock()
	defer storageSync.Unlock()
	val := s.collection[Name]
	val.Name = Name
	val.Type = types.CounterMetricTypeName
	val.ValueInt = val.ValueInt + Value
	s.collection[Name] = val
}
