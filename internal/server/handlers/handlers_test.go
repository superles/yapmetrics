package handlers

import (
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainPage(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "positive test #1",
			want: want{
				code:        200,
				response:    `{"status":"ok"}`,
				contentType: "text/html; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			MainPage(w, request, httprouter.Params{})

			res := w.Result()
			bodyCloseError := res.Body.Close()
			require.NoError(t, bodyCloseError)
			// проверяем код ответа
			assert.Equal(t, res.StatusCode, test.want.code)
			// получаем и проверяем тело запроса

			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			//assert.JSONEq(t, string(resBody), test.want.response)
			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)
		})
	}
}

func TestValuePage(t *testing.T) {
	type want struct {
		code        int
		url         string
		response    string
		contentType string
		value       interface{}
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "value negative test #1",
			want: want{
				code:        404,
				url:         "/value/counter/testSetGet226",
				response:    ``,
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "value negative test #2",
			want: want{
				code:        400,
				url:         "/value/unknown/testSetGet226",
				response:    ``,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.want.url, nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			ValuePage(w, request, httprouter.Params{})

			res := w.Result()
			bodyCloseError := res.Body.Close()
			require.NoError(t, bodyCloseError)
			// проверяем код ответа
			//fmt.Println(test.want, res)
			assert.Equal(t, test.want.code, res.StatusCode, "error code")
			// получаем и проверяем тело запроса

			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, string(resBody), test.want.response)
			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)
		})
	}
}

func TestUnknownPage(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "negative test #1",
			want: want{
				code:        400,
				response:    ``,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/unknown", nil)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			UnknownPage(w, request, httprouter.Params{})

			res := w.Result()
			bodyCloseError := res.Body.Close()
			require.NoError(t, bodyCloseError)
			// проверяем код ответа
			assert.Equal(t, res.StatusCode, test.want.code)
			// получаем и проверяем тело запроса

			_, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			//assert.JSONEq(t, string(resBody), test.want.response)
			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)
		})
	}
}
