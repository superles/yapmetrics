package handlers

import (
	"errors"
	"github.com/superles/yapmetrics/internal/storage"
	"strconv"
)

func counter(name string, value string) (int64, error) {
	// If no name was given, return an error with a message.
	if name == "" {
		return int64(0), errors.New("имя не должно быть пустым")
	}

	v, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		return v, errors.New("значение должно быть int64")
	}

	oldValue := int64(0)

	metric, mErr := storage.MetricRepository.Get(name)

	if mErr == nil {
		oldValue = metric.Value.(int64)
	}

	storage.MetricRepository.Set(name, "counter", v+oldValue)

	return v, nil

}
