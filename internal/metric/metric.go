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
		return strconv.FormatFloat(m.Value, 'f', -1, 64), nil
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
		return "", errors.New("тип метрики не существует")
	}
}

func FromJson(data *JsonData) (*Metric, error) {
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
		return model, errors.New("тип метрики не существует")
	}
	return model, nil
}

func (m *Metric) ToJson() (*JsonData, error) {
	model := &JsonData{ID: m.Name}
	switch m.Type {
	case GaugeMetricType:
		model.MType = GaugeMetricTypeName
		model.Value = &m.Value
	case CounterMetricType:
		model.MType = CounterMetricTypeName
		val := int64(m.Value)
		model.Delta = &val
	default:
		return model, errors.New("тип метрики не существует")
	}
	return model, nil
}
