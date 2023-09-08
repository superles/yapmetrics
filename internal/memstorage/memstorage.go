package memstorage

import (
	"errors"
)

type Metric struct {
	Name  string
	Type  string
	Value string
}

type MetricCollection map[string]Metric

type MemStorage struct {
	Collection MetricCollection
}

func (m *MemStorage) Add(doc Metric) {
	m.Collection[doc.Name] = doc
}

func (m *MemStorage) Get(name string) (Metric, error) {
	val, ok := m.Collection[name]
	if !ok {
		return val, errors.New("Такая метрика отсутствует")
	}
	return m.Collection[name], nil
}

var Storage MemStorage

func init() {
	Storage = MemStorage{Collection: MetricCollection{}}
}
