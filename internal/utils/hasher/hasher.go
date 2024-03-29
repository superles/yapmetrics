package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Encode кодирует массив байтов с использованием hmac sha256 и ключа secret.
func Encode(data []byte, secret []byte) string {
	// Создаем hmac sha256 с ключом
	h := hmac.New(sha256.New, secret)
	h.Write(data)
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func Decode(data []byte, secret []byte) string {
	// Создаем hmac sha256 с ключом
	h := hmac.New(sha256.New, secret)
	h.Write(data)
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
