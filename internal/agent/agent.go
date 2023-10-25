package agent

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/mailru/easyjson"
	"github.com/superles/yapmetrics/internal/agent/client"
	"github.com/superles/yapmetrics/internal/agent/config"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"go.uber.org/zap"
	"io"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type metricProvider interface {
	GetAll(ctx context.Context) (map[string]types.Metric, error)
	SetFloat(ctx context.Context, Name string, Value float64) error
	IncCounter(ctx context.Context, Name string, Value int64) error
}

const attempts = 4

type Agent struct {
	storage metricProvider
	config  *config.Config
	client  client.Client
	logger  *zap.SugaredLogger
}

func New(s metricProvider, cfg *config.Config) *Agent {
	agent := &Agent{storage: s, config: cfg, client: client.NewHTTPAgentClient()}
	return agent
}

func (a *Agent) sendWithRetry(url string, contentType string, body []byte) error {

	start := time.Now()

	response, err := retry.DoWithData(
		func() (*http.Response, error) {
			resp, err := a.client.Post(url, contentType, body, true)

			if err != nil {
				if resp == nil {
					return nil, err
				} else {
					return nil, retry.Unrecoverable(err)
				}
			}

			return resp, nil
		},
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			delay := int(n*2 + 1)
			return time.Duration(delay) * time.Second
		}),
		retry.Attempts(uint(attempts)),
	)

	if response == nil || err != nil {
		return err
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			logger.Log.Error(err)
		}
	}()

	finish := time.Since(start)

	var bodyStr string
	var statusCode int

	bodyReader := response.Body

	contentEncoding := response.Header.Get("Content-Encoding")
	sendsGzip := strings.Contains(contentEncoding, "gzip")

	if sendsGzip {
		zr, err := gzip.NewReader(response.Body)
		if err != nil {
			return err
		}
		bodyReader = zr
	}

	bodyBytes, readErr := io.ReadAll(bodyReader)
	if readErr != nil {
		return readErr
	}
	bodyStr = string(bodyBytes)
	statusCode = response.StatusCode

	logger.Log.Debug("send request",
		zap.String("url", url),
		zap.String("body", string(body)),
		zap.Duration("duration", finish),
		zap.Int("responseCode", statusCode),
		zap.String("responseBody", bodyStr),
		zap.Error(err),
	)

	if response.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("сервер вернул неожиданный код ответа %d", response.StatusCode)
}

func (a *Agent) send(url string, contentType string, body []byte) error {

	start := time.Now()

	response, err := a.client.Post(url, contentType, body, true)

	if response == nil || err != nil {
		return err
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			logger.Log.Error(err)
		}
	}()

	finish := time.Since(start)

	var bodyStr string
	var statusCode int

	bodyReader := response.Body

	contentEncoding := response.Header.Get("Content-Encoding")
	sendsGzip := strings.Contains(contentEncoding, "gzip")

	if sendsGzip {
		zr, err := gzip.NewReader(response.Body)
		if err != nil {
			return err
		}
		bodyReader = zr
	}

	bodyBytes, readErr := io.ReadAll(bodyReader)
	if readErr != nil {
		return readErr
	}
	bodyStr = string(bodyBytes)
	statusCode = response.StatusCode

	logger.Log.Info("send request",
		zap.String("url", url),
		zap.String("body", string(body)),
		zap.Duration("duration", finish),
		zap.Int("responseCode", statusCode),
		zap.String("responseBody", bodyStr),
		zap.Error(err),
	)

	if response.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("сервер вернул неожиданный код ответа %d", response.StatusCode)
}

func (a *Agent) capture() {
	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)
	ctx := context.Background()
	checkError := func(err error) {
		if err != nil {
			logger.Log.Error(err)
		}
	}
	err := a.storage.SetFloat(ctx, "Alloc", float64(stats.Alloc))
	checkError(err)
	err = a.storage.SetFloat(ctx, "BuckHashSys", float64(stats.BuckHashSys))
	checkError(err)
	err = a.storage.SetFloat(ctx, "Frees", float64(stats.Frees))
	checkError(err)
	err = a.storage.SetFloat(ctx, "GCCPUFraction", stats.GCCPUFraction)
	checkError(err)
	err = a.storage.SetFloat(ctx, "GCSys", float64(stats.GCSys))
	checkError(err)
	err = a.storage.SetFloat(ctx, "HeapAlloc", float64(stats.HeapAlloc))
	checkError(err)
	err = a.storage.SetFloat(ctx, "HeapIdle", float64(stats.HeapIdle))
	checkError(err)
	err = a.storage.SetFloat(ctx, "HeapInuse", float64(stats.HeapInuse))
	checkError(err)
	err = a.storage.SetFloat(ctx, "HeapObjects", float64(stats.HeapObjects))
	checkError(err)
	err = a.storage.SetFloat(ctx, "HeapReleased", float64(stats.HeapReleased))
	checkError(err)
	err = a.storage.SetFloat(ctx, "HeapSys", float64(stats.HeapSys))
	checkError(err)
	err = a.storage.SetFloat(ctx, "LastGC", float64(stats.LastGC))
	checkError(err)
	err = a.storage.SetFloat(ctx, "Lookups", float64(stats.Lookups))
	checkError(err)
	err = a.storage.SetFloat(ctx, "MCacheInuse", float64(stats.MCacheInuse))
	checkError(err)
	err = a.storage.SetFloat(ctx, "MCacheSys", float64(stats.MCacheSys))
	checkError(err)
	err = a.storage.SetFloat(ctx, "MSpanInuse", float64(stats.MSpanInuse))
	checkError(err)
	err = a.storage.SetFloat(ctx, "MSpanSys", float64(stats.MSpanSys))
	checkError(err)
	err = a.storage.SetFloat(ctx, "Mallocs", float64(stats.Mallocs))
	checkError(err)
	err = a.storage.SetFloat(ctx, "NextGC", float64(stats.NextGC))
	checkError(err)
	err = a.storage.SetFloat(ctx, "NumForcedGC", float64(stats.NumForcedGC))
	checkError(err)
	err = a.storage.SetFloat(ctx, "NumGC", float64(stats.NumGC))
	checkError(err)
	err = a.storage.SetFloat(ctx, "OtherSys", float64(stats.OtherSys))
	checkError(err)
	err = a.storage.SetFloat(ctx, "PauseTotalNs", float64(stats.PauseTotalNs))
	checkError(err)
	err = a.storage.SetFloat(ctx, "StackInuse", float64(stats.StackInuse))
	checkError(err)
	err = a.storage.SetFloat(ctx, "StackSys", float64(stats.StackSys))
	checkError(err)
	err = a.storage.SetFloat(ctx, "Sys", float64(stats.Sys))
	checkError(err)
	err = a.storage.SetFloat(ctx, "TotalAlloc", float64(stats.TotalAlloc))
	checkError(err)
	err = a.storage.SetFloat(ctx, "RandomValue", 1000+rand.Float64()*(1000-0))
	checkError(err)
	err = a.storage.IncCounter(ctx, "PollCount", 1)
	checkError(err)
}

func (a *Agent) sendPlain(data *types.Metric) error {
	typeStr, _ := types.TypeToString(data.Type)
	strVal, errVal := data.String()
	if errVal != nil {
		return errVal
	}
	url := "http://" + a.config.Endpoint + "/update/" + typeStr + "/" + data.Name + "/" + strVal + ""
	err := a.send(url, "text/plain", []byte(""))
	return err
}

func (a *Agent) sendJSON(data *types.Metric) error {
	updatedJSON, err := data.ToJSON()
	if err != nil {
		return err
	}
	rawBytes, _ := easyjson.Marshal(updatedJSON)
	url := "http://" + a.config.Endpoint + "/update/"
	sendErr := a.sendWithRetry(url, "application/json", rawBytes)
	return sendErr
}

func (a *Agent) sendAllJSON() error {

	logger.Log.Debug("sendAllJSON")

	metrics, err := a.storage.GetAll(context.Background())

	if err != nil {
		return err
	}

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
	sendErr := a.sendWithRetry(url, "application/json", rawBytes)
	return sendErr
}

func (a *Agent) sendAll() error {

	logger.Log.Debug("sendAll")

	metrics, err := a.storage.GetAll(context.Background())

	if err != nil {
		return err
	}

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

	logger.Log.Debug("agent run")

	a.capture()

	go a.poolTick()

	for range time.Tick(time.Second * time.Duration(a.config.ReportInterval)) {
		logger.Log.Debug("agent run sendAllJSON")
		a.sendAllJSON()
	}
}
