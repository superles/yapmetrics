package agent

import (
	"fmt"
	"github.com/superles/yapmetrics/internal/agent/client"
	"github.com/superles/yapmetrics/internal/types"
	"math/rand"
	"runtime"
	"time"
)

const pollInterval = 2
const reportInterval = 10
const serverAddress = "http://localhost:8080"

var metrics types.Metrics

func capture() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	metrics.Alloc = types.Gauge(stats.Alloc)
	metrics.BuckHashSys = types.Gauge(stats.BuckHashSys)
	metrics.Frees = types.Gauge(stats.Frees)
	metrics.GCCPUFraction = types.Gauge(stats.GCCPUFraction)
	metrics.GCSys = types.Gauge(stats.GCSys)
	metrics.HeapAlloc = types.Gauge(stats.HeapAlloc)
	metrics.HeapIdle = types.Gauge(stats.HeapIdle)
	metrics.HeapInuse = types.Gauge(stats.HeapInuse)
	metrics.HeapObjects = types.Gauge(stats.HeapObjects)
	metrics.HeapReleased = types.Gauge(stats.HeapReleased)
	metrics.HeapSys = types.Gauge(stats.HeapSys)
	metrics.LastGC = types.Gauge(stats.LastGC)
	metrics.Lookups = types.Gauge(stats.Lookups)
	metrics.MCacheInuse = types.Gauge(stats.MCacheInuse)
	metrics.MCacheSys = types.Gauge(stats.MCacheSys)
	metrics.MSpanInuse = types.Gauge(stats.MSpanInuse)
	metrics.MSpanSys = types.Gauge(stats.MSpanSys)
	metrics.Mallocs = types.Gauge(stats.Mallocs)
	metrics.NextGC = types.Gauge(stats.NextGC)
	metrics.NumForcedGC = types.Gauge(stats.NumForcedGC)
	metrics.NumGC = types.Gauge(stats.NumGC)
	metrics.OtherSys = types.Gauge(stats.OtherSys)
	metrics.PauseTotalNs = types.Gauge(stats.PauseTotalNs)
	metrics.StackInuse = types.Gauge(stats.StackInuse)
	metrics.StackSys = types.Gauge(stats.StackSys)
	metrics.Sys = types.Gauge(stats.Sys)
	metrics.TotalAlloc = types.Gauge(stats.TotalAlloc)
	metrics.PollCount = metrics.PollCount + 1
	metrics.RandomValue = types.Gauge(1000 + rand.Float64()*(1000-0))
}

func send(mName string, mType string, mValue string) {
	url := serverAddress + "/update/" + mName + "/" + mType + "/" + mValue + ""
	_, err := client.Send(url)
	if err != nil {
		return
	}
}

func sendAll() {
	send("Alloc", "gauge", fmt.Sprintf("%f", metrics.Alloc))
	send("BuckHashSys", "gauge", fmt.Sprintf("%f", metrics.BuckHashSys))
	send("Frees", "gauge", fmt.Sprintf("%f", metrics.Frees))
	send("GCCPUFraction", "gauge", fmt.Sprintf("%f", metrics.GCCPUFraction))
	send("GCSys", "gauge", fmt.Sprintf("%f", metrics.GCSys))
	send("HeapAlloc", "gauge", fmt.Sprintf("%f", metrics.HeapAlloc))
	send("HeapIdle", "gauge", fmt.Sprintf("%f", metrics.HeapIdle))
	send("HeapInuse", "gauge", fmt.Sprintf("%f", metrics.HeapInuse))
	send("HeapObjects", "gauge", fmt.Sprintf("%f", metrics.HeapObjects))
	send("HeapReleased", "gauge", fmt.Sprintf("%f", metrics.HeapReleased))
	send("HeapSys", "gauge", fmt.Sprintf("%f", metrics.HeapSys))
	send("LastGC", "gauge", fmt.Sprintf("%f", metrics.LastGC))
	send("Lookups", "gauge", fmt.Sprintf("%f", metrics.Lookups))
	send("MCacheInuse", "gauge", fmt.Sprintf("%f", metrics.MCacheInuse))
	send("MCacheSys", "gauge", fmt.Sprintf("%f", metrics.MCacheSys))
	send("MSpanInuse", "gauge", fmt.Sprintf("%f", metrics.MSpanInuse))
	send("MSpanSys", "gauge", fmt.Sprintf("%f", metrics.MSpanSys))
	send("Mallocs", "gauge", fmt.Sprintf("%f", metrics.Mallocs))
	send("NextGC", "gauge", fmt.Sprintf("%f", metrics.NextGC))
	send("NumForcedGC", "gauge", fmt.Sprintf("%f", metrics.NumForcedGC))
	send("NumGC", "gauge", fmt.Sprintf("%f", metrics.NumGC))
	send("OtherSys", "gauge", fmt.Sprintf("%f", metrics.OtherSys))
	send("PauseTotalNs", "gauge", fmt.Sprintf("%f", metrics.PauseTotalNs))
	send("StackInuse", "gauge", fmt.Sprintf("%f", metrics.StackInuse))
	send("StackSys", "gauge", fmt.Sprintf("%f", metrics.StackSys))
	send("Sys", "gauge", fmt.Sprintf("%f", metrics.Sys))
	send("TotalAlloc", "gauge", fmt.Sprintf("%f", metrics.TotalAlloc))
	send("TotalAlloc", "gauge", fmt.Sprintf("%f", metrics.TotalAlloc))
	send("TotalAlloc", "gauge", fmt.Sprintf("%f", metrics.TotalAlloc))
	send("RandomValue", "gauge", fmt.Sprintf("%f", metrics.RandomValue))
	send("PollCount", "counter", fmt.Sprintf("%d", metrics.PollCount))
}

func init() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	metrics = types.Metrics{
		Alloc:         types.Gauge(stats.Alloc),
		BuckHashSys:   types.Gauge(stats.BuckHashSys),
		Frees:         types.Gauge(stats.Frees),
		GCCPUFraction: types.Gauge(stats.GCCPUFraction),
		GCSys:         types.Gauge(stats.GCSys),
		HeapAlloc:     types.Gauge(stats.HeapAlloc),
		HeapIdle:      types.Gauge(stats.HeapIdle),
		HeapInuse:     types.Gauge(stats.HeapInuse),
		HeapObjects:   types.Gauge(stats.HeapObjects),
		HeapReleased:  types.Gauge(stats.HeapReleased),
		HeapSys:       types.Gauge(stats.HeapSys),
		LastGC:        types.Gauge(stats.LastGC),
		Lookups:       types.Gauge(stats.Lookups),
		MCacheInuse:   types.Gauge(stats.MCacheInuse),
		MCacheSys:     types.Gauge(stats.MCacheSys),
		MSpanInuse:    types.Gauge(stats.MSpanInuse),
		MSpanSys:      types.Gauge(stats.MSpanSys),
		Mallocs:       types.Gauge(stats.Mallocs),
		NextGC:        types.Gauge(stats.NextGC),
		NumForcedGC:   types.Gauge(stats.NumForcedGC),
		NumGC:         types.Gauge(stats.NumGC),
		OtherSys:      types.Gauge(stats.OtherSys),
		PauseTotalNs:  types.Gauge(stats.PauseTotalNs),
		StackInuse:    types.Gauge(stats.StackInuse),
		StackSys:      types.Gauge(stats.StackSys),
		Sys:           types.Gauge(stats.Sys),
		TotalAlloc:    types.Gauge(stats.TotalAlloc),
		PollCount:     types.Counter(1),
		RandomValue:   types.Gauge(1),
	}
}

func poolTick() {
	for range time.Tick(time.Second * pollInterval) {
		fmt.Println("capture")

		capture()
	}
}

func reportTick() {
	for range time.Tick(time.Second * reportInterval) {
		sendAll()
	}
}

func Run() {

	go poolTick()
	go reportTick()
	//client.Send("http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>")

	time.Sleep(time.Second * reportInterval * 10)
}
