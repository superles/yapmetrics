package memstorage

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"io"
	"os"
	"sync"
)

var storageSync = sync.RWMutex{}

var dumpSync = sync.Mutex{}

type MemStorage struct {
	collection map[string]metric.Metric
}

func New() *MemStorage {
	return &MemStorage{make(map[string]metric.Metric)}
}

func (s *MemStorage) GetAll(ctx context.Context) (map[string]metric.Metric, error) {
	storageSync.RLock()
	defer storageSync.RUnlock()
	targetMap := make(map[string]metric.Metric, len(s.collection))

	// Copy from the original map to the target map
	for key, value := range s.collection {
		targetMap[key] = value
	}

	return targetMap, nil
}

func (s *MemStorage) Get(ctx context.Context, name string) (metric.Metric, error) {
	storageSync.RLock()
	defer storageSync.RUnlock()
	val, ok := s.collection[name]
	if !ok {
		return metric.Metric{}, errors.New("метрика не найдена")
	}
	return val, nil
}

func (s *MemStorage) Set(ctx context.Context, data metric.Metric) error {
	storageSync.Lock()
	defer storageSync.Unlock()
	s.collection[data.Name] = data
	return nil
}

func (s *MemStorage) SetAll(ctx context.Context, data []metric.Metric) error {
	for _, value := range data {
		logger.Log.Debug("SetAll", value)
		switch value.Type {
		case metric.GaugeMetricType:
			if err := s.SetFloat(ctx, value.Name, value.Value); err != nil {
				return err
			}
		case metric.CounterMetricType:
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
	s.collection[Name] = metric.Metric{Name: Name, Type: metric.GaugeMetricType, Value: Value}
	return nil
}

func (s *MemStorage) IncCounter(ctx context.Context, Name string, Value int64) error {
	storageSync.Lock()
	defer storageSync.Unlock()
	val := s.collection[Name]
	val.Name = Name
	val.Type = metric.CounterMetricType
	val.Value = val.Value + float64(Value)
	s.collection[Name] = val
	return nil
}

func (s *MemStorage) Ping(ctx context.Context) error {
	return nil
}

func (s *MemStorage) Dump(ctx context.Context, path string) error {
	dumpSync.Lock()
	defer dumpSync.Unlock()
	file, fileErr := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0666)
	if fileErr != nil {
		return fileErr
	}
	//if err := file.Truncate(0); err != nil {
	//	return err
	//}
	encoder := json.NewEncoder(file)
	metrics, err := s.GetAll(ctx)
	if err != nil {
		return err
	}
	for _, item := range metrics {
		err := encoder.Encode(&item)
		if err != nil {
			return fileErr
		}
	}

	if err := file.Sync(); err != nil {
		return err
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
		var m metric.Metric
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		if err := s.Set(context.Background(), m); err != nil {
			return err
		}
	}

	if err := file.Close(); err != nil {
		return err
	}
	logger.Log.Info("load success")
	return nil
}
