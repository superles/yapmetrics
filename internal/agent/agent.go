package agent

import (
	"fmt"
	"github.com/superles/yapmetrics/internal/agent/client"
	"github.com/superles/yapmetrics/internal/agent/config"
	types "github.com/superles/yapmetrics/internal/metric"
	"math/rand"
	"runtime"
	"time"
)

type metricProvider interface {
	GetAll() map[string]types.Metric
	SetFloat(Name string, Value float64)
	IncCounter(Name string, Value int64)
}

type Agent struct {
	storage metricProvider
	config  *config.Config
}

func New(s metricProvider) *Agent {
	agent := &Agent{storage: s, config: config.New()}
	return agent
}

func (a *Agent) capture() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	a.storage.SetFloat("Alloc", float64(stats.Alloc))
	a.storage.SetFloat("BuckHashSys", float64(stats.BuckHashSys))
	a.storage.SetFloat("Frees", float64(stats.Frees))
	a.storage.SetFloat("GCCPUFraction", stats.GCCPUFraction)
	a.storage.SetFloat("GCSys", float64(stats.GCSys))
	a.storage.SetFloat("HeapAlloc", float64(stats.HeapAlloc))
	a.storage.SetFloat("HeapIdle", float64(stats.HeapIdle))
	a.storage.SetFloat("HeapInuse", float64(stats.HeapInuse))
	a.storage.SetFloat("HeapObjects", float64(stats.HeapObjects))
	a.storage.SetFloat("HeapReleased", float64(stats.HeapReleased))
	a.storage.SetFloat("HeapSys", float64(stats.HeapSys))
	a.storage.SetFloat("LastGC", float64(stats.LastGC))
	a.storage.SetFloat("Lookups", float64(stats.Lookups))
	a.storage.SetFloat("MCacheInuse", float64(stats.MCacheInuse))
	a.storage.SetFloat("MCacheSys", float64(stats.MCacheSys))
	a.storage.SetFloat("MSpanInuse", float64(stats.MSpanInuse))
	a.storage.SetFloat("MSpanSys", float64(stats.MSpanSys))
	a.storage.SetFloat("Mallocs", float64(stats.Mallocs))
	a.storage.SetFloat("NextGC", float64(stats.NextGC))
	a.storage.SetFloat("NumForcedGC", float64(stats.NumForcedGC))
	a.storage.SetFloat("NumGC", float64(stats.NumGC))
	a.storage.SetFloat("OtherSys", float64(stats.OtherSys))
	a.storage.SetFloat("PauseTotalNs", float64(stats.PauseTotalNs))
	a.storage.SetFloat("StackInuse", float64(stats.StackInuse))
	a.storage.SetFloat("StackSys", float64(stats.StackSys))
	a.storage.SetFloat("Sys", float64(stats.Sys))
	a.storage.SetFloat("TotalAlloc", float64(stats.TotalAlloc))
	a.storage.SetFloat("RandomValue", 1000+rand.Float64()*(1000-0))
	a.storage.IncCounter("PollCount", 1)
}

func (a *Agent) send(mName string, mType string, mValue string) error {
	url := "http://" + a.config.Endpoint + "/update/" + mType + "/" + mName + "/" + mValue + ""
	_, err := client.Send(url)
	return err
}

func (a *Agent) sendAll() error {
	fmt.Println("sendAll")

	metrics := a.storage.GetAll()

	for Name, Item := range metrics {
		strVal, errVal := Item.String()
		if errVal != nil {
			return errVal
		}
		err := a.send(Name, Item.Type, strVal)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Agent) poolTick() {
	for range time.Tick(time.Duration(a.config.PollInterval) * time.Second) {
		fmt.Println("capture")
		a.capture()
	}
}

func (a *Agent) Run() {

	a.capture()

	go a.poolTick()

	for range time.Tick(time.Second * time.Duration(a.config.ReportInterval)) {
		a.sendAll()
	}
}
