package client

import (
	"errors"
	"github.com/gojek/heimdall/v7/httpclient"
	"log"
	"net/http"
)

const useHttp = true

func postHeimdall(url string) (*http.Response, error) {
	var client = httpclient.NewClient()
	headers := http.Header{}
	headers.Set("Content-Type", "text/plain")
	return client.Post(url, nil, headers)
}

func postHTTP(url string) (*http.Response, error) {
	client := http.Client{}
	return client.Post(url, "text/plain", nil)
}

func Send(url string) (bool, error) {
	post := postHeimdall
	if useHttp {
		post = postHTTP
	}
	response, postErr := post(url)
	if postErr != nil {
		return false, postErr
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	if response.StatusCode == http.StatusOK {
		return true, nil
	}

	return false, errors.New("unknown error")
}
