package client

import (
	"bytes"
	"net/http"
)

type HTTPAgentClient struct {
	client *http.Client
}

func (c *HTTPAgentClient) Post(url string, contentType string, body []byte) (*http.Response, error) {
	return c.client.Post(url, contentType, bytes.NewReader(body))
}

func NewHTTPAgentClient() Client {
	return &HTTPAgentClient{&http.Client{}}
}
