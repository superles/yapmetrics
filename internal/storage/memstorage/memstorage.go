package memstorage

import (
	"context"
	"encoding/json"
	"errors"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"io"
	"os"
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

func (s *MemStorage) Dump(ctx context.Context, path string) error {
	file, fileErr := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if fileErr != nil {
		return fileErr
	}
	if err := file.Truncate(0); err != nil {
		return err
	}
	encoder := json.NewEncoder(file)
	metrics, err := s.GetAll(ctx)
	if err != nil {
		return err
	}
	for _, metric := range metrics {
		err := encoder.Encode(&metric)
		if err != nil {
			return fileErr
		}
	}
	if err := file.Close(); err != nil {
		return err
	}
	logger.Log.Info("dump success")
	return nil
}

func (s *MemStorage) Restore(ctx context.Context, path string) error {
	file, fileErr := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0666)
	if fileErr != nil {
		return fileErr
	}
	dec := json.NewDecoder(file)
	for {
		var m types.Metric
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if err := s.Set(context.Background(), &m); err != nil {
			return err
		}
	}

	if err := file.Close(); err != nil {
		return err
	}
	logger.Log.Info("load success")
	return nil
}
