package agent

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"github.com/avast/retry-go/v4"
	"github.com/mailru/easyjson"
	"github.com/superles/yapmetrics/internal/agent/client"
	"github.com/superles/yapmetrics/internal/agent/config"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const attempts = 4

type Agent struct {
	storage metricProvider
	config  *config.Config
	client  client.Client
	logger  *zap.SugaredLogger
}

// New Создание нового агента.
func New(s metricProvider, cfg *config.Config) *Agent {
	agent := &Agent{storage: s, config: cfg, client: client.NewHTTPAgentClient(client.AgentClientParams{Key: cfg.SecretKey})}
	return agent
}

func (a *Agent) sendWithRetry(url string, contentType string, body []byte) error {

	start := time.Now()

	response, err := retry.DoWithData(
		func() (*http.Response, error) {
			resp, err := a.client.Post(url, contentType, body, true)
			if err != nil {
				var opError *net.OpError
				if errors.As(err, &opError) {
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

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("сервер вернул неожиданный код ответа %d", response.StatusCode)
	}

	return nil
}

func (a *Agent) sendPlain(data *types.Metric) error {
	typeStr, _ := types.TypeToString(data.Type)
	strVal, errVal := data.String()
	if errVal != nil {
		return errVal
	}
	url := fmt.Sprintf("http://%s/update/%s/%s/%s", a.config.Endpoint, typeStr, data.Name, strVal)
	err := a.send(url, "text/plain", []byte(""))
	return err
}

func (a *Agent) sendJSON(data *types.Metric) error {
	updatedJSON, err := data.ToJSON()
	if err != nil {
		return err
	}
	rawBytes, _ := easyjson.Marshal(updatedJSON)
	url := fmt.Sprintf("http://%s/update/", a.config.Endpoint)
	sendErr := a.sendWithRetry(url, "application/json", rawBytes)
	return sendErr
}

func (a *Agent) sendAllJSON(ctx context.Context) error {

	logger.Log.Debug("sendAllJSON")

	metrics, err := a.storage.GetAll(ctx)

	if err != nil {
		return err
	}

	rawBytes, err := compressMetrics(metrics)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://%s/updates/", a.config.Endpoint)
	err = a.sendWithRetry(url, "application/json", rawBytes)
	return err
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

func (a *Agent) poolTickRuntime(ctx context.Context, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logger.Log.Debug("captureRuntime")
			err := a.captureRuntime(ctx)
			if err != nil {
				logger.Log.Error("captureRuntime error", err.Error())
			}
		}
	}
}

func (a *Agent) poolTickPsutil(ctx context.Context, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logger.Log.Debug("capturePsutil")
			err := a.capturePsutil(ctx)
			if err != nil {
				logger.Log.Error("capturePsutil error", err.Error())
			}
		}
	}
}

func (a *Agent) sendTicker(ctx context.Context, reportInterval time.Duration) {
	ticker := time.NewTicker(reportInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			logger.Log.Debug("agent run sendAllJSON")
			err := a.sendAllJSON(ctx)
			if err != nil {
				logger.Log.Error("reportInterval error", err.Error())
			}
		}
	}
}

// Run Запуск агента.
func (a *Agent) Run(ctx context.Context) error {

	logger.Log.Debug("agent run")

	reportInterval := time.Second * time.Duration(a.config.ReportInterval)
	pollInterval := time.Second * time.Duration(a.config.PollInterval)

	go a.poolTickRuntime(ctx, pollInterval)
	go a.poolTickPsutil(ctx, pollInterval)

	if a.config.RateLimit > 0 {
		a.sendPoolTicker(ctx, reportInterval)
	} else {
		go a.sendTicker(ctx, reportInterval)
	}

	return nil
}
