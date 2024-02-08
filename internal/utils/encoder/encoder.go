package encoder

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

type Encoder struct {
	publicKey *rsa.PublicKey
}

func New(publicKeyPath string) (*Encoder, error) {
	publicKeyFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	// Извлечение открытого ключа из PEM-блока
	block, _ := pem.Decode(publicKeyFile)
	if block == nil {
		return nil, errors.New("ошибка извлечения ключа из pem")
	}

	// Преобразование открытого ключа в структуру rsa.PublicKey
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return &Encoder{publicKey: publicKey}, nil
}

func minCustom(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (e Encoder) Encrypt(messageBytes []byte) ([]byte, error) {
	// Определение максимальной длины блока для RSA шифрования
	maxBlockSize := e.publicKey.N.BitLen()/8 - 11

	var encryptedBytes []byte

	// Шифрование блоков сообщения
	for len(messageBytes) > 0 {
		blockSize := minCustom(len(messageBytes), maxBlockSize)
		block, err := rsa.EncryptPKCS1v15(rand.Reader, e.publicKey, messageBytes[:blockSize])
		if err != nil {
			return []byte(""), err
		}
		encryptedBytes = append(encryptedBytes, block...)
		messageBytes = messageBytes[blockSize:]
	}

	return encryptedBytes, nil
}
