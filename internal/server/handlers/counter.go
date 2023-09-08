package handlers

import (
	"errors"
	"fmt"
	"github.com/superles/yapmetrics/internal/memstorage"
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

	metric, mErr := memstorage.Storage.Get(name)

	if mErr != nil {
		memstorage.Storage.Add(memstorage.Metric{
			Name:  name,
			Type:  "counter",
			Value: value,
		})
	} else {
		metricValue, metricValueErr := strconv.ParseInt(metric.Value, 10, 64)
		if metricValueErr != nil {
			metricValue = 0
		}
		memstorage.Storage.Add(memstorage.Metric{
			Name:  name,
			Type:  "counter",
			Value: fmt.Sprint(v + metricValue),
		})
	}

	return v, nil

}
