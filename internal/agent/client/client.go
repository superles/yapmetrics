package client

import (
	"errors"
	"log"
	"net/http"
)

//func postHeimdall(url string) (*http.Response, error) {
//	client := httpclient.NewClient()
//	headers := http.Header{}
//	headers.Set("Content-Type", "text/plain")
//	return client.Post(url, nil, headers)
//}

func postHTTP(url string) (*http.Response, error) {
	client := http.Client{}
	return client.Post(url, "text/plain", nil)
}

func Send(url string) (bool, error) {
	response, postErr := postHTTP(url)
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
