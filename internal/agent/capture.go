package agent

import (
	"context"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/superles/yapmetrics/internal/metric"
	"math/rand"
	"runtime"
	"strconv"
)

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
func generateRandomValue() float64 {
	return randFloat(0, 1000)
}

func (a *Agent) captureRuntime(ctx context.Context) error {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	metrics := make([]metric.Metric, 0)
	metrics = append(metrics, metric.Metric{Name: "Alloc", Type: metric.GaugeMetricType, Value: float64(stats.Alloc)})
	metrics = append(metrics, metric.Metric{Name: "BuckHashSys", Type: metric.GaugeMetricType, Value: float64(stats.BuckHashSys)})
	metrics = append(metrics, metric.Metric{Name: "Frees", Type: metric.GaugeMetricType, Value: float64(stats.Frees)})
	metrics = append(metrics, metric.Metric{Name: "GCCPUFraction", Type: metric.GaugeMetricType, Value: stats.GCCPUFraction})
	metrics = append(metrics, metric.Metric{Name: "GCSys", Type: metric.GaugeMetricType, Value: float64(stats.GCSys)})
	metrics = append(metrics, metric.Metric{Name: "HeapAlloc", Type: metric.GaugeMetricType, Value: float64(stats.HeapAlloc)})
	metrics = append(metrics, metric.Metric{Name: "HeapIdle", Type: metric.GaugeMetricType, Value: float64(stats.HeapIdle)})
	metrics = append(metrics, metric.Metric{Name: "HeapInuse", Type: metric.GaugeMetricType, Value: float64(stats.HeapInuse)})
	metrics = append(metrics, metric.Metric{Name: "HeapObjects", Type: metric.GaugeMetricType, Value: float64(stats.HeapObjects)})
	metrics = append(metrics, metric.Metric{Name: "HeapReleased", Type: metric.GaugeMetricType, Value: float64(stats.HeapReleased)})
	metrics = append(metrics, metric.Metric{Name: "HeapSys", Type: metric.GaugeMetricType, Value: float64(stats.HeapSys)})
	metrics = append(metrics, metric.Metric{Name: "LastGC", Type: metric.GaugeMetricType, Value: float64(stats.LastGC)})
	metrics = append(metrics, metric.Metric{Name: "Lookups", Type: metric.GaugeMetricType, Value: float64(stats.Lookups)})
	metrics = append(metrics, metric.Metric{Name: "MCacheInuse", Type: metric.GaugeMetricType, Value: float64(stats.MCacheInuse)})
	metrics = append(metrics, metric.Metric{Name: "MCacheSys", Type: metric.GaugeMetricType, Value: float64(stats.MCacheSys)})
	metrics = append(metrics, metric.Metric{Name: "MSpanInuse", Type: metric.GaugeMetricType, Value: float64(stats.MSpanInuse)})
	metrics = append(metrics, metric.Metric{Name: "MSpanSys", Type: metric.GaugeMetricType, Value: float64(stats.MSpanSys)})
	metrics = append(metrics, metric.Metric{Name: "Mallocs", Type: metric.GaugeMetricType, Value: float64(stats.Mallocs)})
	metrics = append(metrics, metric.Metric{Name: "NextGC", Type: metric.GaugeMetricType, Value: float64(stats.NextGC)})
	metrics = append(metrics, metric.Metric{Name: "NumForcedGC", Type: metric.GaugeMetricType, Value: float64(stats.NumForcedGC)})
	metrics = append(metrics, metric.Metric{Name: "NumGC", Type: metric.GaugeMetricType, Value: float64(stats.NumGC)})
	metrics = append(metrics, metric.Metric{Name: "OtherSys", Type: metric.GaugeMetricType, Value: float64(stats.OtherSys)})
	metrics = append(metrics, metric.Metric{Name: "PauseTotalNs", Type: metric.GaugeMetricType, Value: float64(stats.PauseTotalNs)})
	metrics = append(metrics, metric.Metric{Name: "StackInuse", Type: metric.GaugeMetricType, Value: float64(stats.StackInuse)})
	metrics = append(metrics, metric.Metric{Name: "StackSys", Type: metric.GaugeMetricType, Value: float64(stats.StackSys)})
	metrics = append(metrics, metric.Metric{Name: "Sys", Type: metric.GaugeMetricType, Value: float64(stats.Sys)})
	metrics = append(metrics, metric.Metric{Name: "TotalAlloc", Type: metric.GaugeMetricType, Value: float64(stats.TotalAlloc)})
	metrics = append(metrics, metric.Metric{Name: "RandomValue", Type: metric.GaugeMetricType, Value: generateRandomValue()})
	metrics = append(metrics, metric.Metric{Name: "RandomValue", Type: metric.CounterMetricType, Value: 1})

	err := a.storage.SetAll(ctx, metrics)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) capturePsutil(ctx context.Context) error {

	metrics := make([]metric.Metric, 0)

	stats, err := mem.VirtualMemory()
	if err != nil {
		return err
	}

	percents, err := cpu.Percent(0, true)
	if err != nil {
		return err
	}

	metrics = append(metrics, metric.Metric{Name: "TotalMemory", Type: metric.GaugeMetricType, Value: float64(stats.Total)})
	metrics = append(metrics, metric.Metric{Name: "FreeMemory", Type: metric.GaugeMetricType, Value: float64(stats.Free)})
	for idx, percent := range percents {
		name := "CPUutilization" + strconv.Itoa(idx+1)
		metrics = append(metrics, metric.Metric{Name: name, Type: metric.GaugeMetricType, Value: percent})
	}

	err = a.storage.SetAll(ctx, metrics)

	if err != nil {
		return err
	}
	return nil
}
