package handlers

import (
	"errors"
	"github.com/superles/yapmetrics/internal/storage"
	"strconv"
)

func gauge(name string, value string) (float64, error) {
	// If no name was given, return an error with a message.
	if name == "" {
		return float64(0), errors.New("имя не должно быть пустым")
	}

	v, err := strconv.ParseFloat(value, 64)

	if err != nil {
		return v, errors.New("значение должно быть float64")
	}

	storage.MetricRepository.Set(name, "gauge", v)

	return v, nil

}
