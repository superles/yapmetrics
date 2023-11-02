package agent

import (
	"context"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"math/rand"
	"runtime"
	"strconv"
)

func (a *Agent) captureRuntime(ctx context.Context) error {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	err := a.storage.SetFloat(ctx, "Alloc", float64(stats.Alloc))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "BuckHashSys", float64(stats.BuckHashSys))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "Frees", float64(stats.Frees))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "GCCPUFraction", stats.GCCPUFraction)
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "GCSys", float64(stats.GCSys))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "HeapAlloc", float64(stats.HeapAlloc))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "HeapIdle", float64(stats.HeapIdle))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "HeapInuse", float64(stats.HeapInuse))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "HeapObjects", float64(stats.HeapObjects))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "HeapReleased", float64(stats.HeapReleased))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "HeapSys", float64(stats.HeapSys))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "LastGC", float64(stats.LastGC))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "Lookups", float64(stats.Lookups))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "MCacheInuse", float64(stats.MCacheInuse))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "MCacheSys", float64(stats.MCacheSys))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "MSpanInuse", float64(stats.MSpanInuse))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "MSpanSys", float64(stats.MSpanSys))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "Mallocs", float64(stats.Mallocs))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "NextGC", float64(stats.NextGC))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "NumForcedGC", float64(stats.NumForcedGC))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "NumGC", float64(stats.NumGC))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "OtherSys", float64(stats.OtherSys))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "PauseTotalNs", float64(stats.PauseTotalNs))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "StackInuse", float64(stats.StackInuse))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "StackSys", float64(stats.StackSys))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "Sys", float64(stats.Sys))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "TotalAlloc", float64(stats.TotalAlloc))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "RandomValue", 1000+rand.Float64()*(1000-0))
	if err != nil {
		return err
	}
	err = a.storage.IncCounter(ctx, "PollCount", 1)
	if err != nil {
		return err
	}

	return nil
}

func (a *Agent) capturePsutil(ctx context.Context) error {
	stats, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "TotalMemory", float64(stats.Total))
	if err != nil {
		return err
	}
	err = a.storage.SetFloat(ctx, "FreeMemory", float64(stats.Free))
	if err != nil {
		return err
	}

	percents, err := cpu.Percent(0, true)

	if err != nil {
		return err
	}
	for idx, percent := range percents {
		name := "CPUutilization" + strconv.Itoa(idx+1)
		err = a.storage.SetFloat(ctx, name, percent)
		if err != nil {
			return err
		}
	}
	return nil
}
