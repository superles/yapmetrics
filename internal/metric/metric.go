package metric

import (
	"fmt"
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

type Collection map[string]Metric

func (m *Metric) String() (string, error) {
	switch m.Type {
	case GaugeMetricType:
		return strconv.FormatFloat(m.Value, 'f', -1, 64), nil
	case CounterMetricType:
		return strconv.FormatFloat(m.Value, 'f', 0, 64), nil
	default:
		return "", fmt.Errorf("ошибка вывода значения метрики: %d", m.Type)
	}
}

func StringToType(mType string) (int, error) {
	switch mType {
	case GaugeMetricTypeName:
		return GaugeMetricType, nil
	case CounterMetricTypeName:
		return CounterMetricType, nil
	default:
		return 0, fmt.Errorf("ошибка вывода значения метрики: %s", mType)
	}
}

func TypeToString(mType int) (string, error) {
	switch mType {
	case GaugeMetricType:
		return GaugeMetricTypeName, nil
	case CounterMetricType:
		return CounterMetricTypeName, nil
	default:
		return "", fmt.Errorf("тип метрики не существует: %d", mType)
	}
}

func FromJSON(data *JSONData) (*Metric, error) {
	model := &Metric{Name: data.ID}
	switch data.MType {
	case GaugeMetricTypeName:
		model.Type = GaugeMetricType
		if data.Value != nil {
			model.Value = *data.Value
		}
		model.Value = *data.Value
	case CounterMetricTypeName:
		model.Type = CounterMetricType
		if data.Delta != nil {
			model.Value = float64(*data.Delta)
		}

	default:
		return model, fmt.Errorf("тип метрики не существует: %s", data.MType)
	}
	return model, nil
}

func (m *Metric) ToJSON() (*JSONData, error) {
	model := &JSONData{ID: m.Name}
	switch m.Type {
	case GaugeMetricType:
		model.MType = GaugeMetricTypeName
		val := m.Value
		model.Value = &val
	case CounterMetricType:
		model.MType = CounterMetricTypeName
		val := int64(m.Value)
		model.Delta = &val
	default:
		return model, fmt.Errorf("тип метрики не существует: %d", m.Type)
	}
	return model, nil
}
