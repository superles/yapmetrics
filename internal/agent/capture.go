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

func (a *Agent) captureRuntime(ctx context.Context) error {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	metrics := []metric.Metric{
		{Name: "Alloc", Type: metric.GaugeMetricType, Value: float64(stats.Alloc)},
		{Name: "BuckHashSys", Type: metric.GaugeMetricType, Value: float64(stats.BuckHashSys)},
		{Name: "Frees", Type: metric.GaugeMetricType, Value: float64(stats.Frees)},
		{Name: "GCCPUFraction", Type: metric.GaugeMetricType, Value: stats.GCCPUFraction},
		{Name: "GCSys", Type: metric.GaugeMetricType, Value: float64(stats.GCSys)},
		{Name: "HeapAlloc", Type: metric.GaugeMetricType, Value: float64(stats.HeapAlloc)},
		{Name: "HeapIdle", Type: metric.GaugeMetricType, Value: float64(stats.HeapIdle)},
		{Name: "HeapInuse", Type: metric.GaugeMetricType, Value: float64(stats.HeapInuse)},
		{Name: "HeapObjects", Type: metric.GaugeMetricType, Value: float64(stats.HeapObjects)},
		{Name: "HeapReleased", Type: metric.GaugeMetricType, Value: float64(stats.HeapReleased)},
		{Name: "HeapSys", Type: metric.GaugeMetricType, Value: float64(stats.HeapSys)},
		{Name: "LastGC", Type: metric.GaugeMetricType, Value: float64(stats.LastGC)},
		{Name: "Lookups", Type: metric.GaugeMetricType, Value: float64(stats.Lookups)},
		{Name: "MCacheInuse", Type: metric.GaugeMetricType, Value: float64(stats.MCacheInuse)},
		{Name: "MCacheSys", Type: metric.GaugeMetricType, Value: float64(stats.MCacheSys)},
		{Name: "MSpanInuse", Type: metric.GaugeMetricType, Value: float64(stats.MSpanInuse)},
		{Name: "MSpanSys", Type: metric.GaugeMetricType, Value: float64(stats.MSpanSys)},
		{Name: "Mallocs", Type: metric.GaugeMetricType, Value: float64(stats.Mallocs)},
		{Name: "NextGC", Type: metric.GaugeMetricType, Value: float64(stats.NextGC)},
		{Name: "NumForcedGC", Type: metric.GaugeMetricType, Value: float64(stats.NumForcedGC)},
		{Name: "NumGC", Type: metric.GaugeMetricType, Value: float64(stats.NumGC)},
		{Name: "OtherSys", Type: metric.GaugeMetricType, Value: float64(stats.OtherSys)},
		{Name: "PauseTotalNs", Type: metric.GaugeMetricType, Value: float64(stats.PauseTotalNs)},
		{Name: "StackInuse", Type: metric.GaugeMetricType, Value: float64(stats.StackInuse)},
		{Name: "StackSys", Type: metric.GaugeMetricType, Value: float64(stats.StackSys)},
		{Name: "Sys", Type: metric.GaugeMetricType, Value: float64(stats.Sys)},
		{Name: "TotalAlloc", Type: metric.GaugeMetricType, Value: float64(stats.TotalAlloc)},
		{Name: "RandomValue", Type: metric.GaugeMetricType, Value: randFloat(0, 1000)},
		{Name: "PollCount", Type: metric.CounterMetricType, Value: 1},
	}

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
