package agent

import (
	"fmt"
	"github.com/superles/yapmetrics/internal/agent/client"
	"github.com/superles/yapmetrics/internal/agent/config"
	"github.com/superles/yapmetrics/internal/storage"
	"github.com/superles/yapmetrics/internal/types"
	"math/rand"
	"runtime"
	"strings"
	"time"
)

var pollCount = 0

func formatFloat(gauge interface{}) string {
	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%f", gauge), "0"), ".")
}

func capture(count int) {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	storage.MetricRepository.Set("Alloc", "gauge", types.Gauge(stats.Alloc))
	storage.MetricRepository.Set("BuckHashSys", "gauge", types.Gauge(stats.BuckHashSys))
	storage.MetricRepository.Set("Frees", "gauge", types.Gauge(stats.Frees))
	storage.MetricRepository.Set("GCCPUFraction", "gauge", types.Gauge(stats.GCCPUFraction))
	storage.MetricRepository.Set("GCSys", "gauge", types.Gauge(stats.GCSys))
	storage.MetricRepository.Set("HeapAlloc", "gauge", types.Gauge(stats.HeapAlloc))
	storage.MetricRepository.Set("HeapIdle", "gauge", types.Gauge(stats.HeapIdle))
	storage.MetricRepository.Set("HeapInuse", "gauge", types.Gauge(stats.HeapInuse))
	storage.MetricRepository.Set("HeapObjects", "gauge", types.Gauge(stats.HeapObjects))
	storage.MetricRepository.Set("HeapReleased", "gauge", types.Gauge(stats.HeapReleased))
	storage.MetricRepository.Set("HeapSys", "gauge", types.Gauge(stats.HeapSys))
	storage.MetricRepository.Set("LastGC", "gauge", types.Gauge(stats.LastGC))
	storage.MetricRepository.Set("Lookups", "gauge", types.Gauge(stats.Lookups))
	storage.MetricRepository.Set("MCacheInuse", "gauge", types.Gauge(stats.MCacheInuse))
	storage.MetricRepository.Set("MCacheSys", "gauge", types.Gauge(stats.MCacheSys))
	storage.MetricRepository.Set("MSpanInuse", "gauge", types.Gauge(stats.MSpanInuse))
	storage.MetricRepository.Set("MSpanSys", "gauge", types.Gauge(stats.MSpanSys))
	storage.MetricRepository.Set("Mallocs", "gauge", types.Gauge(stats.Mallocs))
	storage.MetricRepository.Set("NextGC", "gauge", types.Gauge(stats.NextGC))
	storage.MetricRepository.Set("NumForcedGC", "gauge", types.Gauge(stats.NumForcedGC))
	storage.MetricRepository.Set("NumGC", "gauge", types.Gauge(stats.NumGC))
	storage.MetricRepository.Set("OtherSys", "gauge", types.Gauge(stats.OtherSys))
	storage.MetricRepository.Set("PauseTotalNs", "gauge", types.Gauge(stats.PauseTotalNs))
	storage.MetricRepository.Set("StackInuse", "gauge", types.Gauge(stats.StackInuse))
	storage.MetricRepository.Set("StackSys", "gauge", types.Gauge(stats.StackSys))
	storage.MetricRepository.Set("Sys", "gauge", types.Gauge(stats.Sys))
	storage.MetricRepository.Set("TotalAlloc", "gauge", types.Gauge(stats.TotalAlloc))
	storage.MetricRepository.Set("RandomValue", "gauge", types.Gauge(1000+rand.Float64()*(1000-0)))
	storage.MetricRepository.Set("PollCount", "counter", types.Counter(count))
}

func send(mName string, mType string, mValue string) error {
	url := "http://" + config.AgentConfig.Endpoint + "/update/" + mType + "/" + mName + "/" + mValue + ""
	_, err := client.Send(url)
	return err
}

func sendAll() error {
	fmt.Println("sendAll")

	metrics := storage.MetricRepository.GetAll()

	for Name, Item := range metrics {
		if Item.Type == "counter" {
			err := send(Name, Item.Type, fmt.Sprintf("%d", Item.Value))
			if err != nil {
				return err
			}
		} else {
			err := send(Name, Item.Type, formatFloat(Item.Value))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func init() {
	capture(0)
}

func poolTick() {
	for range time.Tick(time.Duration(config.AgentConfig.PollInterval) * time.Second) {
		fmt.Println("capture")
		pollCount = pollCount + 1
		capture(pollCount)
	}
}

func Run() {

	config.InitConfig()

	go poolTick()

	for range time.Tick(time.Second * time.Duration(config.AgentConfig.ReportInterval)) {
		sendAll()
	}
}
