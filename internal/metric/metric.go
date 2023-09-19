package metric

import (
	"errors"
	"strconv"
)

const (
	GaugeMetricTypeName   = "gauge"
	CounterMetricTypeName = "counter"
)

const (
	GaugeMetricType   = 1
	CounterMetricType = 2
)

type Gauge float64
type Counter int64

type Metric struct {
	Name  string  //имя метрики
	Type  int     //тип метрики counter | gauge
	Value float64 //Значение метрики
}

func (m *Metric) String() (string, error) {
	switch m.Type {
	case GaugeMetricType:
		return strconv.FormatFloat(m.Value, 'g', -1, 64), nil
	case CounterMetricType:
		return strconv.FormatFloat(m.Value, 'f', 0, 64), nil
	default:
		return "", errors.New("ошибка вывода значения метрики")
	}
}

func StringToType(mType string) (int, error) {
	switch mType {
	case GaugeMetricTypeName:
		return GaugeMetricType, nil
	case CounterMetricTypeName:
		return CounterMetricType, nil
	default:
		return 0, errors.New("ошибка вывода значения метрики")
	}
}

func TypeToString(mType int) (string, error) {
	switch mType {
	case GaugeMetricType:
		return GaugeMetricTypeName, nil
	case CounterMetricType:
		return CounterMetricTypeName, nil
	default:
		return "", errors.New("ошибка вывода значения метрики")
	}
}
