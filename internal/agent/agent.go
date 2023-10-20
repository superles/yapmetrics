package agent

import (
	"compress/gzip"
	"context"
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
	"strings"
	"time"
)

type metricProvider interface {
	GetAll(ctx context.Context) map[string]types.Metric
	SetFloat(ctx context.Context, Name string, Value float64)
	IncCounter(ctx context.Context, Name string, Value int64)
}

type Agent struct {
	storage metricProvider
	config  *config.Config
	client  client.Client
}

func New(s metricProvider, cfg *config.Config) *Agent {
	err := logger.Initialize(cfg.LogLevel)
	if err != nil {
		log.Panicln("ошибка инициализации логера", err.Error())
	}
	agent := &Agent{storage: s, config: cfg, client: client.NewHTTPAgentClient()}
	return agent
}

func (a *Agent) send(url string, contentType string, body []byte) (bool, error) {

	start := time.Now()

	response, postErr := a.client.Post(url, contentType, body, true)

	finish := time.Since(start)

	var bodyStr string
	var statusCode int
	if response != nil {

		bodyReader := response.Body

		contentEncoding := response.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")

		if sendsGzip {
			zr, err := gzip.NewReader(response.Body)
			if err != nil {
				return false, err
			}
			bodyReader = zr
		}

		bodyBytes, readErr := io.ReadAll(bodyReader)
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
	ctx := context.Background()
	a.storage.SetFloat(ctx, "Alloc", float64(stats.Alloc))
	a.storage.SetFloat(ctx, "BuckHashSys", float64(stats.BuckHashSys))
	a.storage.SetFloat(ctx, "Frees", float64(stats.Frees))
	a.storage.SetFloat(ctx, "GCCPUFraction", stats.GCCPUFraction)
	a.storage.SetFloat(ctx, "GCSys", float64(stats.GCSys))
	a.storage.SetFloat(ctx, "HeapAlloc", float64(stats.HeapAlloc))
	a.storage.SetFloat(ctx, "HeapIdle", float64(stats.HeapIdle))
	a.storage.SetFloat(ctx, "HeapInuse", float64(stats.HeapInuse))
	a.storage.SetFloat(ctx, "HeapObjects", float64(stats.HeapObjects))
	a.storage.SetFloat(ctx, "HeapReleased", float64(stats.HeapReleased))
	a.storage.SetFloat(ctx, "HeapSys", float64(stats.HeapSys))
	a.storage.SetFloat(ctx, "LastGC", float64(stats.LastGC))
	a.storage.SetFloat(ctx, "Lookups", float64(stats.Lookups))
	a.storage.SetFloat(ctx, "MCacheInuse", float64(stats.MCacheInuse))
	a.storage.SetFloat(ctx, "MCacheSys", float64(stats.MCacheSys))
	a.storage.SetFloat(ctx, "MSpanInuse", float64(stats.MSpanInuse))
	a.storage.SetFloat(ctx, "MSpanSys", float64(stats.MSpanSys))
	a.storage.SetFloat(ctx, "Mallocs", float64(stats.Mallocs))
	a.storage.SetFloat(ctx, "NextGC", float64(stats.NextGC))
	a.storage.SetFloat(ctx, "NumForcedGC", float64(stats.NumForcedGC))
	a.storage.SetFloat(ctx, "NumGC", float64(stats.NumGC))
	a.storage.SetFloat(ctx, "OtherSys", float64(stats.OtherSys))
	a.storage.SetFloat(ctx, "PauseTotalNs", float64(stats.PauseTotalNs))
	a.storage.SetFloat(ctx, "StackInuse", float64(stats.StackInuse))
	a.storage.SetFloat(ctx, "StackSys", float64(stats.StackSys))
	a.storage.SetFloat(ctx, "Sys", float64(stats.Sys))
	a.storage.SetFloat(ctx, "TotalAlloc", float64(stats.TotalAlloc))
	a.storage.SetFloat(ctx, "RandomValue", 1000+rand.Float64()*(1000-0))
	a.storage.IncCounter(ctx, "PollCount", 1)
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
	updatedJSON, err := data.ToJSON()
	if err != nil {
		return err
	}
	rawBytes, _ := easyjson.Marshal(updatedJSON)
	url := "http://" + a.config.Endpoint + "/update/"
	_, sendErr := a.send(url, "application/json", rawBytes)
	return sendErr
}

func (a *Agent) sendAllJSON() error {

	logger.Log.Debug("sendAllJSON")

	metrics := a.storage.GetAll(context.Background())

	var col types.JSONDataCollection
	for _, item := range metrics {
		updatedJSON, err := item.ToJSON()
		if err != nil {
			return err
		}
		col = append(col, *updatedJSON)
	}
	rawBytes, _ := easyjson.Marshal(col)
	url := "http://" + a.config.Endpoint + "/updates/"
	_, sendErr := a.send(url, "application/json", rawBytes)
	return sendErr
}

func (a *Agent) sendAll() error {

	logger.Log.Debug("sendAll")

	metrics := a.storage.GetAll(context.Background())

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

	logger.Log.Sugar().Debug("agent run")

	a.capture()

	go a.poolTick()

	for range time.Tick(time.Second * time.Duration(a.config.ReportInterval)) {
		logger.Log.Sugar().Debug("agent run sendAllJSON")
		err := a.sendAllJSON()
		if err != nil {
			return
		}
	}
}
