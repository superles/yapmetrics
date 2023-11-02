package agent

import (
	"context"
	"errors"
	"github.com/mailru/easyjson"
	types "github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"sync"
	"time"
)

type response struct {
	Status   int
	Error    error
	WorkerId int
}

func (a *Agent) generator(ctx context.Context) <-chan types.Collection {
	ch := make(chan types.Collection, 1)
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(a.config.ReportInterval))
		defer func() {
			ticker.Stop()
			close(ch)
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

				select {
				case <-ch:
					logger.Log.Debug("generator free chanel")
				default:
				}

				ch <- metrics
			}
		}
	}()
	return ch
}

func (a *Agent) worker(id int, ctx context.Context, input <-chan types.Collection, results chan<- response) {
Loop:
	for {
		select {
		case <-ctx.Done():
			return // Выход из горутины при отмене контекста
		case metrics, ok := <-input:

			if !ok {
				results <- response{WorkerId: id, Error: errors.New("input channel closed")}
				return
			}

			var col types.JSONDataCollection
			for _, item := range metrics {
				updatedJSON, err := item.ToJSON()
				if err != nil {
					results <- response{WorkerId: id, Error: err}
					logger.Log.Debug("loop continue")
					continue Loop
				}
				col = append(col, *updatedJSON)
			}
			rawBytes, err := easyjson.Marshal(col)
			if err != nil {
				results <- response{WorkerId: id, Error: err}
				logger.Log.Debug("loop continue")
				continue Loop
			}
			url := "http://" + a.config.Endpoint + "/updates/"
			sendErr := a.sendWithRetry(url, "application/json", rawBytes)
			if sendErr != nil {
				logger.Log.Debug("loop continue")
				results <- response{WorkerId: id, Error: sendErr}
			}
		}
	}
}

func (a *Agent) sendPoolTicker(ctx context.Context) {
	var wg sync.WaitGroup
	requestChan := a.generator(ctx)
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
		close(resultChan)
	}()
	go func() {
		for resp := range resultChan {
			if resp.Error != nil {
				logger.Log.Error("worker error", resp.WorkerId, resp.Error.Error())
			}
		}
	}()
}
