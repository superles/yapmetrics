package metric

import (
	"errors"
	"fmt"
)

const (
	GaugeMetricTypeName   = "gauge"
	CounterMetricTypeName = "counter"
)

type Gauge float64
type Counter int64

type Metric struct {
	Name       string //имя метрики
	Type       string //тип метрики counter | gauge
	ValueInt   int64  //Значение метрики
	ValueFloat float64
}

func (m *Metric) String() (string, error) {
	switch m.Type {
	case GaugeMetricTypeName:
		return fmt.Sprintf("%g", m.ValueFloat), nil
	case CounterMetricTypeName:
		return fmt.Sprintf("%d", m.ValueInt), nil
	default:
		return "", errors.New("ошибка вывода значения метрики")
	}
}
