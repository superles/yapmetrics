package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func Encode(data string, secret string) string {

	// Создаем hmac sha256 с ключом
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}

func Decode(data []byte, secret string) string {

	// Создаем hmac sha256 с ключом
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	return sha
}
