package server

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/server/config"
	"github.com/superles/yapmetrics/internal/storage"
	"github.com/superles/yapmetrics/internal/storage/memstorage"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_GetValue(t *testing.T) {

	repo := memstorage.New()
	cfg := config.New()
	serv := New(repo, cfg)
	ts := httptest.NewServer(serv.router)
	defer ts.Close()

	type want struct {
		code int
	}
	tests := []struct {
		name    string
		method  string
		url     string
		storage storage.Storage
		want    want
	}{
		{
			name:   "test #1",
			method: http.MethodGet,
			url:    "/value/counter/testSetGet226",
			want: want{
				code: 400,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, ts.URL+test.url, nil)
			require.NoError(t, err)

			resp, err := ts.Client().Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()
		})
	}
}

func TestServer_UpdateCounter(t *testing.T) {

	repo := memstorage.New()
	cfg := config.New()
	serv := New(repo, cfg)
	ts := httptest.NewServer(serv.router)
	defer ts.Close()

	type want struct {
		code  int
		value float64
	}
	tests := []struct {
		name    string
		method  string
		url     string
		storage storage.Storage
		want    want
	}{
		{
			name:   "test #1",
			method: http.MethodGet,
			url:    "/update/counter/testSetGet226/22",
			want: want{
				code:  400,
				value: 22,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, ts.URL+test.url, nil)
			require.NoError(t, err)
			resp, err := ts.Client().Do(req)
			result, _ := repo.Get(context.Background(), "testSetGet226")
			if test.want.value != result.Value {
				t.Error("значение в хранилище не соответсвует")
			}
			require.NoError(t, err)
			defer resp.Body.Close()
		})
	}
}

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
