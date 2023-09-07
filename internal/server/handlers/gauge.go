package handlers

import (
	"errors"
	"strconv"
)

func gauge(name string, value string) (float64, error) {
	// If no name was given, return an error with a message.
	if name == "" {
		return float64(0), errors.New("имя не должно быть пустым")
	}

	if v, err := strconv.ParseFloat(value, 64); err == nil {
		return v, nil
	} else {
		return v, errors.New("значение должно быть float64")
	}

}
