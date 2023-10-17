package memstorage

import (
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

func (s *MemStorage) Get(name string) (types.Metric, bool) {
	storageSync.Lock()
	defer storageSync.Unlock()
	val, ok := s.collection[name]
	if !ok {
		return types.Metric{}, false
	}
	return val, true
}

func (s *MemStorage) Set(data *types.Metric) {
	storageSync.Lock()
	defer storageSync.Unlock()
	s.collection[data.Name] = *data
}

func (s *MemStorage) SetFloat(Name string, Value float64) {
	storageSync.Lock()
	defer storageSync.Unlock()
	s.collection[Name] = types.Metric{Name: Name, Type: types.GaugeMetricType, Value: Value}
}

func (s *MemStorage) IncCounter(Name string, Value int64) {
	storageSync.Lock()
	defer storageSync.Unlock()
	val := s.collection[Name]
	val.Name = Name
	val.Type = types.CounterMetricType
	val.Value = val.Value + float64(Value)
	s.collection[Name] = val
}
