package decrypt

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/superles/yapmetrics/internal/utils/logger"
	"io"
	"net/http"
	"os"
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

func minCustom(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func decryptMessage(encryptedBytes []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	// Определение максимальной длины блока для RSA дешифрования
	maxBlockSize := privateKey.N.BitLen() / 8

	// Дешифрование блоков сообщения
	var decryptedBytes []byte
	for len(encryptedBytes) > 0 {
		blockSize := minCustom(len(encryptedBytes), maxBlockSize)
		decryptedBlock, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, encryptedBytes[:blockSize])
		if err != nil {
			return decryptedBytes, err
		}
		decryptedBytes = append(decryptedBytes, decryptedBlock...)
		encryptedBytes = encryptedBytes[blockSize:]
	}

	return decryptedBytes, nil
}

func WithDecrypt(secretFile string) func(h http.Handler) http.Handler {
	defaultHandler := func(h http.Handler) http.Handler {
		return h
	}
	if len(secretFile) == 0 {
		return defaultHandler
	}
	privateKeyFile, err := os.ReadFile(secretFile)
	if err != nil {
		logger.Log.Error("Ошибка чтения закрытого ключа:", err)
		return defaultHandler
	}
	// Извлечение закрытого ключа из PEM-блока
	block, _ := pem.Decode(privateKeyFile)
	if block == nil {
		logger.Log.Error("Ошибка декодирования PEM-блока закрытого ключа")
		return defaultHandler
	}
	// Преобразование закрытого ключа в структуру rsa.PrivateKey
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		logger.Log.Error("Ошибка парсинга закрытого ключа:", err)
		return defaultHandler
	}

	return func(h http.Handler) http.Handler {
		logFn := func(w http.ResponseWriter, r *http.Request) {
			// Чтение тела запроса
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Ошибка чтения тела запроса", http.StatusBadRequest)
				return
			}

			// Расшифровка тела запроса
			decryptedBody, err := decryptMessage(body, privateKey)
			if err != nil {
				http.Error(w, "Ошибка расшифровки запроса", http.StatusBadRequest)
				return
			}

			// Замена тела запроса на расшифрованные данные
			r.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))

			// Передача управления следующему обработчику
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(logFn)
	}
}
