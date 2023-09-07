package client

import (
	"errors"
	"net/http"
)

func Send(url string) (bool, error) {
	response, _ := http.Post(url, "text/plain", nil)
	err := response.Body.Close()
	if err != nil {
		return false, err
	}

	if response.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, errors.New("unknown error")
}
