package client

import (
	"bytes"
	"net/http"
)

type HttpAgentClient struct {
	client *http.Client
}

func (c *HttpAgentClient) Post(url string, contentType string, body []byte) (*http.Response, error) {
	return c.client.Post(url, contentType, bytes.NewReader(body))
}

func NewHttpAgentClient() Client {
	return &HttpAgentClient{&http.Client{}}
}
