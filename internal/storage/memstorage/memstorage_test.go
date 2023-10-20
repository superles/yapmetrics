package memstorage

import (
	"context"
	types "github.com/superles/yapmetrics/internal/metric"
	"reflect"
	"testing"
)

func TestMemStorage_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		initMap map[string]types.Metric
		want    map[string]types.Metric
		isEqual bool
	}{
		{
			"test empty positive #1",
			map[string]types.Metric{},
			map[string]types.Metric{},
			false,
		},
		{
			"test empty negative #2",
			map[string]types.Metric{
				"test": {
					Name:  "test",
					Type:  types.CounterMetricType,
					Value: float64(0),
				},
			},
			map[string]types.Metric{},
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := &MemStorage{
				collection: test.initMap,
			}
			if got := s.GetAll(context.Background()); reflect.DeepEqual(got, test.want) == test.isEqual {
				t.Errorf("GetAll() = %v, want %v", got, test.want)
			}
		})
	}
}
