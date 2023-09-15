package server

import (
	"github.com/stretchr/testify/require"
	"github.com/superles/yapmetrics/internal/storage"
	"github.com/superles/yapmetrics/internal/storage/memstorage"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_GetValue(t *testing.T) {

	repo := memstorage.New()
	serv := New(repo)
	ts := httptest.NewServer(serv.Router)
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
	serv := New(repo)
	ts := httptest.NewServer(serv.Router)
	defer ts.Close()

	type want struct {
		code     int
		valueInt int64
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
				code:     400,
				valueInt: 22,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest(test.method, ts.URL+test.url, nil)
			require.NoError(t, err)
			resp, err := ts.Client().Do(req)
			result, _ := repo.Get("testSetGet226")
			if test.want.valueInt != result.ValueInt {
				t.Error("значение в хранилище не соответсвует")
			}
			require.NoError(t, err)
			defer resp.Body.Close()
		})
	}
}
