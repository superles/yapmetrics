package logging

import (
	"bytes"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type (
	// берём структуру для хранения сведений об ответе
	responseData struct {
		status      int
		size        int
		body        string
		contentType string
	}

	// добавляем реализацию http.ResponseWriter
	loggingResponseWriter struct {
		http.ResponseWriter // встраиваем оригинальный http.ResponseWriter
		responseData        *responseData
	}
)

func (r *loggingResponseWriter) Header() http.Header {
	return r.ResponseWriter.Header()
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	// записываем ответ, используя оригинальный http.ResponseWriter
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size // захватываем размер
	r.responseData.body = string(b)
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	// записываем код статуса, используя оригинальный http.ResponseWriter
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // захватываем код статуса
	r.responseData.contentType = r.ResponseWriter.Header().Get("Content-Type")
}

func WithLogging(h http.Handler) http.Handler {

	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w, // встраиваем оригинальный http.ResponseWriter
			responseData:   responseData,
		}
		var bodyStr string

		if r.Method == http.MethodPost {
			bodyBytes, _ := io.ReadAll(r.Body)
			err := r.Body.Close()
			if err != nil {
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			bodyStr = string(bodyBytes)
		}

		h.ServeHTTP(&lw, r) // внедряем реализацию http.ResponseWriter

		duration := time.Since(start)

		logger.Log.Debug("receive request",
			zap.String("uri", r.RequestURI),
			zap.String("method", r.Method),
			zap.Int("status", responseData.status), // получаем перехваченный код статуса ответа
			zap.Duration("duration", duration),
			zap.Int("size", responseData.size), // получаем перехваченный размер ответа
			zap.String("requestBody", bodyStr),
			zap.String("responseType", responseData.contentType),
			zap.String("responseBody", responseData.body), //тело ответа если есть
		)
	}
	return http.HandlerFunc(logFn)
}
