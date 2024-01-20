package memstorage

import (
	"context"
	types "github.com/superles/yapmetrics/internal/metric"
	"os"
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
			if got, _ := s.GetAll(context.Background()); reflect.DeepEqual(got, test.want) == test.isEqual {
				t.Errorf("GetAll() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestMemStorage(t *testing.T) {
	// Create a temporary file for testing dump and restore
	tmpfile, err := os.CreateTemp("", "memstorage_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// Initialize MemStorage
	memStorage := New()

	t.Run("Set and Get", func(t *testing.T) {

		testMetric := types.Metric{Name: "test_metric_set_get", Type: types.GaugeMetricType, Value: 42.0}
		err := memStorage.Set(context.Background(), testMetric)

		if err != nil {
			t.Fatalf("Error setting metric: %v", err)
		}
		result, err := memStorage.Get(context.Background(), testMetric.Name)
		if err != nil {
			t.Fatalf("Error getting metric: %v", err)
		}

		if result != testMetric {
			t.Errorf("Expected %v, got %v", testMetric, result)
		}
	})

	t.Run("SetFloat and Get", func(t *testing.T) {

		testMetric := types.Metric{Name: "test_metric_set_float"}

		err := memStorage.SetFloat(context.Background(), testMetric.Name, 50.0)
		if err != nil {
			t.Fatalf("Error setting float metric: %v", err)
		}

		result, err := memStorage.Get(context.Background(), testMetric.Name)
		if err != nil {
			t.Fatalf("Error getting metric: %v", err)
		}

		expected := types.Metric{Name: testMetric.Name, Type: types.GaugeMetricType, Value: 50.0}
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("IncCounter and Get", func(t *testing.T) {

		testCounterMetric := types.Metric{Name: "counter_metric", Type: types.CounterMetricType, Value: 10.0}

		err := memStorage.Set(context.Background(), testCounterMetric)

		if err != nil {
			t.Fatalf("Error setting metric: %v", err)
		}

		err = memStorage.IncCounter(context.Background(), testCounterMetric.Name, 5)
		if err != nil {
			t.Fatalf("Error incrementing counter: %v", err)
		}
		result, err := memStorage.Get(context.Background(), testCounterMetric.Name)
		if err != nil {
			t.Fatalf("Error getting metric: %v", err)
		}

		expected := types.Metric{Name: testCounterMetric.Name, Type: types.CounterMetricType, Value: 15.0}
		if result != expected {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("Dump and Restore", func(t *testing.T) {

		testMetric := types.Metric{Name: "test_metric_dump", Type: types.GaugeMetricType, Value: 42.0}
		err := memStorage.Set(context.Background(), testMetric)
		if err != nil {
			t.Fatalf("Error set metrics: %v", err)
		}
		// Dump metrics to a file
		err = memStorage.Dump(context.Background(), tmpfile.Name())

		if err != nil {
			t.Fatalf("Error dumping metrics: %v", err)
		}

		// Create a new MemStorage for restore
		restoredMemStorage := New()

		// Restore metrics from the file
		err = restoredMemStorage.Restore(context.Background(), tmpfile.Name())
		if err != nil {
			t.Fatalf("Error restoring metrics: %v", err)
		}

		// Verify if the restored metric matches the original
		restoredResult, err := restoredMemStorage.Get(context.Background(), testMetric.Name)
		if err != nil {
			t.Fatalf("Error getting restored metric: %v", err)
		}

		if restoredResult != testMetric {
			t.Errorf("Expected %v, got %v", testMetric, restoredResult)
		}
	})

	// Cleanup: Close and remove the temporary file
	tmpfile.Close()
}
