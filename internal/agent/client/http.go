package client

import (
	"bytes"
	"net/http"
)

type HTTPAgentClient struct {
	client *http.Client
}

func (c *HTTPAgentClient) Post(url string, contentType string, body []byte, compress bool) (*http.Response, error) {
	if !compress {
		return c.client.Post(url, contentType, bytes.NewReader(body))
	}
	cBody, cErr := Compress(body)
	if cErr != nil {
		return nil, cErr
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(cBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	return c.client.Do(req)
}

func NewHTTPAgentClient() Client {
	return &HTTPAgentClient{&http.Client{}}
}
