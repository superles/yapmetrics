package server

import (
	"context"
	"fmt"
	"github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/server/config"
	"github.com/superles/yapmetrics/internal/storage/memstorage"
	"io"
	"net/http"
	"net/http/httptest"
)

func ExampleServer_GetValue() {

	repo := memstorage.New()
	cfg := config.New()
	serv := New(repo, cfg)
	ts := httptest.NewServer(serv.router)
	defer ts.Close()

	err := repo.Set(context.Background(), metric.Metric{
		Name:  "testSetGet226",
		Type:  metric.CounterMetricType,
		Value: 1,
	})

	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodGet, ts.URL+"/value/counter/testSetGet226", nil)

	if err != nil {
		return
	}

	resp, err := ts.Client().Do(req)

	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	// Выводим результаты в соответствии с godoc
	fmt.Println(resp.StatusCode)
	fmt.Println(string(body))
	// Output:
	// 200
	// 1
}
