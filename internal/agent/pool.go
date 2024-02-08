package agent

import (
	"context"
	"errors"
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

func (a *Agent) generator(ctx context.Context, ch chan<- types.Collection, reportInterval time.Duration) {
	ticker := time.NewTicker(reportInterval)
	defer func() {
		ticker.Stop()
		logger.Log.Debug("stop generator ticker and close input channel")
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
			ch <- metrics
		}
	}
}

func (a *Agent) dispatcher(ctx context.Context, input <-chan types.Collection, out chan<- types.Collection) {
	for {
		select {
		case <-ctx.Done():
			return // Выход из горутины при отмене контекста
		case metrics, ok := <-input:
			if !ok {
				logger.Log.Debug("dispatcher input channel closed")
				return
			}

			select {
			case out <- metrics:
			default:
			}
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
			err := a.sendAll(ctx, metrics.ToSlice())
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
	go a.generator(ctx, requestChan, reportInterval)
	dispatcherChan := make(chan types.Collection, a.config.RateLimit)
	go a.dispatcher(ctx, requestChan, dispatcherChan)
	resultChan := make(chan response, a.config.RateLimit)
	for i := 1; i <= int(a.config.RateLimit); i++ {
		go func(workerID int) {
			defer wg.Done()
			a.worker(workerID, ctx, dispatcherChan, resultChan)
		}(i)
	}
	wg.Add(int(a.config.RateLimit))
	go func() {
		wg.Wait()
		logger.Log.Debug("free all channels")
		close(requestChan)
		close(dispatcherChan)
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
