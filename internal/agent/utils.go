package agent

import (
	"github.com/mailru/easyjson"
	types "github.com/superles/yapmetrics/internal/metric"
)

func compressMetrics(metrics types.Collection) ([]byte, error) {
	var col types.JSONDataCollection
	for _, item := range metrics {
		updatedJSON, err := item.ToJSON()
		if err != nil {
			return nil, err
		}
		col = append(col, updatedJSON)
	}
	rawBytes, err := easyjson.Marshal(col)

	if err != nil {
		return nil, err
	}
	return rawBytes, nil
}
