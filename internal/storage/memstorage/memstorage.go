package memstorage

import (
	"context"
	"errors"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"sync"
)

var storageSync = sync.Mutex{}

type MemStorage struct {
	collection map[string]types.Metric
}

func New() *MemStorage {
	return &MemStorage{make(map[string]types.Metric)}
}

func (s *MemStorage) GetAll(ctx context.Context) (map[string]types.Metric, error) {
	storageSync.Lock()
	defer storageSync.Unlock()
	targetMap := make(map[string]types.Metric, len(s.collection))

	// Copy from the original map to the target map
	for key, value := range s.collection {
		targetMap[key] = value
	}

	return targetMap, nil
}

func (s *MemStorage) Get(ctx context.Context, name string) (types.Metric, error) {
	storageSync.Lock()
	defer storageSync.Unlock()
	val, ok := s.collection[name]
	if !ok {
		return types.Metric{}, errors.New("метрика не найдена")
	}
	return val, nil
}

func (s *MemStorage) Set(ctx context.Context, data *types.Metric) error {
	storageSync.Lock()
	defer storageSync.Unlock()
	if data == nil {
		return errors.New("метрика содержит пустой объект")
	}
	s.collection[data.Name] = *data
	return nil
}

func (s *MemStorage) SetAll(ctx context.Context, data *[]types.Metric) error {
	for _, value := range *data {
		logger.Log.Debug("SetAll", value)
		switch value.Type {
		case types.GaugeMetricType:
			if err := s.SetFloat(ctx, value.Name, value.Value); err != nil {
				return err
			}
		case types.CounterMetricType:
			if err := s.IncCounter(ctx, value.Name, int64(value.Value)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *MemStorage) SetFloat(ctx context.Context, Name string, Value float64) error {
	storageSync.Lock()
	defer storageSync.Unlock()
	s.collection[Name] = types.Metric{Name: Name, Type: types.GaugeMetricType, Value: Value}
	return nil
}

func (s *MemStorage) IncCounter(ctx context.Context, Name string, Value int64) error {
	storageSync.Lock()
	defer storageSync.Unlock()
	val := s.collection[Name]
	val.Name = Name
	val.Type = types.CounterMetricType
	val.Value = val.Value + float64(Value)
	s.collection[Name] = val
	return nil
}

func (s *MemStorage) Ping(ctx context.Context) error {
	return nil
}
