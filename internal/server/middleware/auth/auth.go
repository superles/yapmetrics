package auth

import (
	"bytes"
	"fmt"
	"github.com/superles/yapmetrics/internal/utils/hasher"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"io"
	"net/http"
)

type rwWrapper struct {
	http.ResponseWriter
	buf    bytes.Buffer
	status int
	header http.Header
}

func (w *rwWrapper) Write(p []byte) (int, error) {
	return w.buf.Write(p)
}
func (w *rwWrapper) WriteHeader(status int) {
	w.status = status
}

func (w *rwWrapper) Header() http.Header {
	return w.header
}

func WithAuth(key string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		logFn := func(w http.ResponseWriter, r *http.Request) {

			wrapper := rwWrapper{
				ResponseWriter: w,
				status:         200,
				header:         make(http.Header),
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

			hash := r.Header.Get("HashSHA256")
			if len(key) != 0 && len(hash) != 0 {
				expectedHash := hasher.Encode(bodyStr, key)
				if len(hash) == 0 {
					logger.Log.Error("пустой хеш в запросе")
					http.Error(w, "", http.StatusBadRequest)
					return
				} else if hash != expectedHash {
					logger.Log.Error("неправильный хеш в запросе")
					http.Error(w, "", http.StatusBadRequest)
					return
				}
			}

			h.ServeHTTP(&wrapper, r)

			if len(key) != 0 {
				hash := hasher.Encode(string(wrapper.buf.Bytes()), key)
				w.Header().Set("HashSHA256", hash)
			}
			for hKey, hVal := range wrapper.header {
				val := hVal[0]
				w.Header().Set(hKey, val)
			}
			w.WriteHeader(wrapper.status)
			_, err := w.Write(wrapper.buf.Bytes())
			if err != nil {
				logger.Log.Error(fmt.Sprintf("ошибка записи: %s", err.Error()))
			}
		}
		return http.HandlerFunc(logFn)
	}
}
