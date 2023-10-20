package metric

//easyjson:json
type JSONData struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

//easyjson:json
type JSONDataCollection []JSONData

func (v JSONDataCollection) ToMetrics() *[]Metric {

	var targetMap []Metric

	for _, item := range v {
		m := Metric{
			Name: item.ID,
		}
		switch item.MType {
		case GaugeMetricTypeName:
			m.Value = *item.Value
			m.Type = GaugeMetricType
		case CounterMetricTypeName:
			m.Value = float64(*item.Delta)
			m.Type = CounterMetricType
		}
		targetMap = append(targetMap, m)
	}
	return &targetMap
}
