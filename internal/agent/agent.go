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
	Storage metricProvider
	Config  *config.Config
}

func New(s metricProvider) *Agent {
	agent := &Agent{Storage: s, Config: config.New()}
	return agent
}

func (a *Agent) capture() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	a.Storage.SetFloat("Alloc", float64(stats.Alloc))
	a.Storage.SetFloat("BuckHashSys", float64(stats.BuckHashSys))
	a.Storage.SetFloat("Frees", float64(stats.Frees))
	a.Storage.SetFloat("GCCPUFraction", stats.GCCPUFraction)
	a.Storage.SetFloat("GCSys", float64(stats.GCSys))
	a.Storage.SetFloat("HeapAlloc", float64(stats.HeapAlloc))
	a.Storage.SetFloat("HeapIdle", float64(stats.HeapIdle))
	a.Storage.SetFloat("HeapInuse", float64(stats.HeapInuse))
	a.Storage.SetFloat("HeapObjects", float64(stats.HeapObjects))
	a.Storage.SetFloat("HeapReleased", float64(stats.HeapReleased))
	a.Storage.SetFloat("HeapSys", float64(stats.HeapSys))
	a.Storage.SetFloat("LastGC", float64(stats.LastGC))
	a.Storage.SetFloat("Lookups", float64(stats.Lookups))
	a.Storage.SetFloat("MCacheInuse", float64(stats.MCacheInuse))
	a.Storage.SetFloat("MCacheSys", float64(stats.MCacheSys))
	a.Storage.SetFloat("MSpanInuse", float64(stats.MSpanInuse))
	a.Storage.SetFloat("MSpanSys", float64(stats.MSpanSys))
	a.Storage.SetFloat("Mallocs", float64(stats.Mallocs))
	a.Storage.SetFloat("NextGC", float64(stats.NextGC))
	a.Storage.SetFloat("NumForcedGC", float64(stats.NumForcedGC))
	a.Storage.SetFloat("NumGC", float64(stats.NumGC))
	a.Storage.SetFloat("OtherSys", float64(stats.OtherSys))
	a.Storage.SetFloat("PauseTotalNs", float64(stats.PauseTotalNs))
	a.Storage.SetFloat("StackInuse", float64(stats.StackInuse))
	a.Storage.SetFloat("StackSys", float64(stats.StackSys))
	a.Storage.SetFloat("Sys", float64(stats.Sys))
	a.Storage.SetFloat("TotalAlloc", float64(stats.TotalAlloc))
	a.Storage.SetFloat("RandomValue", 1000+rand.Float64()*(1000-0))
	a.Storage.IncCounter("PollCount", 1)
}

func (a *Agent) send(mName string, mType string, mValue string) error {
	url := "http://" + a.Config.Endpoint + "/update/" + mType + "/" + mName + "/" + mValue + ""
	_, err := client.Send(url)
	return err
}

func (a *Agent) sendAll() error {
	fmt.Println("sendAll")

	metrics := a.Storage.GetAll()

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
	for range time.Tick(time.Duration(a.Config.PollInterval) * time.Second) {
		fmt.Println("capture")
		a.capture()
	}
}

func (a *Agent) Run() {

	a.capture()

	go a.poolTick()

	for range time.Tick(time.Second * time.Duration(a.Config.ReportInterval)) {
		a.sendAll()
	}
}
