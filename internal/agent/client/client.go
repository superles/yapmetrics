package client

import "net/http"

type Client interface {
	Post(url string, contentType string, body []byte) (*http.Response, error)
}
