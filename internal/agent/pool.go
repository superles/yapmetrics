package agent

import (
	"context"
	"errors"
	"fmt"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"sync"
	"time"
)

type response struct {
	Status   int
	Error    error
	WorkerID int
}

func (a *Agent) generator(ch chan types.Collection, ctx context.Context, reportInterval time.Duration) {
	ticker := time.NewTicker(reportInterval)
	defer func() {
		ticker.Stop()
		logger.Log.Debug("stop generator ticker")
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			metrics, err := a.storage.GetAll(ctx)
			if err != nil {
				logger.Log.Error("generator GetAll error", err.Error())
				continue
			}

			select {
			case _, ok := <-ch:
				if !ok {
					logger.Log.Error("канал генератора закрыт")
					return
				}
				logger.Log.Debug("generator free chanel")
			default:
			}

			ch <- metrics
		}
	}
}

func (a *Agent) worker(id int, ctx context.Context, input <-chan types.Collection, results chan<- response) {
	for {
		select {
		case <-ctx.Done():
			return // Выход из горутины при отмене контекста
		case metrics, ok := <-input:
			if !ok {
				results <- response{WorkerID: id, Error: errors.New("input channel closed")}
				return
			}

			var rawBytes, err = compressMetrics(metrics)
			if err != nil {
				results <- response{WorkerID: id, Error: err}
				logger.Log.Debug("ошибка сжатия метрик", err.Error())
				continue
			}
			url := fmt.Sprintf("http://%s/updates/", a.config.Endpoint)
			err = a.sendWithRetry(url, "application/json", rawBytes)
			if err != nil {
				logger.Log.Debug("ошибка отправки метрик", err.Error())
				results <- response{WorkerID: id, Error: err}
			}
		}
	}
}

func (a *Agent) sendPoolTicker(ctx context.Context, reportInterval time.Duration) {
	var wg sync.WaitGroup
	requestChan := make(chan types.Collection, 1)
	go a.generator(requestChan, ctx, reportInterval)
	resultChan := make(chan response)
	for i := 1; i <= int(a.config.RateLimit); i++ {
		go func(workerID int) {
			a.worker(workerID, ctx, requestChan, resultChan)
			wg.Done()
		}(i)
	}
	wg.Add(int(a.config.RateLimit))
	go func() {
		wg.Wait()
		logger.Log.Debug("free response channel")
		close(requestChan)
		close(resultChan)
	}()
	go func() {
		for resp := range resultChan {
			if resp.Error != nil {
				logger.Log.Error("worker error", resp.WorkerID, resp.Error.Error())
			}
		}
	}()
}
