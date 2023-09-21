package agent

import (
	"errors"
	"github.com/mailru/easyjson"
	"github.com/superles/yapmetrics/internal/agent/client"
	"github.com/superles/yapmetrics/internal/agent/config"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"go.uber.org/zap"
	"io"
	"log"
	"math/rand"
	"net/http"
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
	client  *client.Client
}

func New(s metricProvider) *Agent {
	cfg := config.New()
	err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Panicln("ошибка инициализации логера", err.Error())
	}
	cl := client.NewHttpAgentClient()
	agent := &Agent{storage: s, config: cfg, client: &cl}
	return agent
}

func (a *Agent) send(url string, contentType string, body []byte) (bool, error) {

	start := time.Now()

	response, postErr := (*a.client).Post(url, contentType, body)

	finish := time.Since(start)

	var bodyStr string
	var statusCode int
	if response != nil {
		bodyBytes, readErr := io.ReadAll(response.Body)
		if readErr != nil {
			logger.Log.Error(readErr.Error())
		}
		bodyStr = string(bodyBytes)
		statusCode = response.StatusCode
	}

	logger.Log.Info("send request",
		zap.String("url", url),
		zap.String("body", string(body)),
		zap.Duration("duration", finish),
		zap.Int("responseCode", statusCode),
		zap.String("responseBody", bodyStr),
		zap.Error(postErr),
	)

	if postErr != nil {
		return false, postErr
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if response.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, errors.New("unknown error")
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

func (a *Agent) sendPlain(data *types.Metric) error {
	typeStr, _ := types.TypeToString(data.Type)
	strVal, errVal := data.String()
	if errVal != nil {
		return errVal
	}
	url := "http://" + a.config.Endpoint + "/update/" + typeStr + "/" + data.Name + "/" + strVal + ""
	_, err := a.send(url, "text/plain", []byte(""))
	return err
}

func (a *Agent) sendJSON(data *types.Metric) error {
	updatedJSON, err := data.ToJson()
	if err != nil {
		return err
	}
	rawBytes, _ := easyjson.Marshal(updatedJSON)
	url := "http://" + a.config.Endpoint + "/update/"
	_, sendErr := a.send(url, "application/json", rawBytes)
	return sendErr
}

func (a *Agent) sendAll() error {

	logger.Log.Debug("sendAll")

	metrics := a.storage.GetAll()

	for _, Item := range metrics {
		err := a.sendJSON(&Item)
		if err != nil {
			logger.Log.Error(err.Error(),
				zap.String("name", Item.Name),
				zap.Int("type", Item.Type),
				zap.Float64("value", Item.Value),
			)
			return err
		}
	}
	return nil
}

func (a *Agent) poolTick() {
	for range time.Tick(time.Duration(a.config.PollInterval) * time.Second) {
		logger.Log.Debug("capture")
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
