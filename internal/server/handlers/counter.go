package handlers

import (
	"errors"
	"strconv"
)

func counter(name string, value string) (int64, error) {
	// If no name was given, return an error with a message.
	if name == "" {
		return int64(0), errors.New("имя не должно быть пустым")
	}

	if v, err := strconv.ParseInt(value, 10, 64); err == nil {
		return v, nil
	} else {
		return v, errors.New("значение должно быть int64")
	}

}
