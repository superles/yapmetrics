package client

import (
	"bytes"
	"github.com/gojek/heimdall/v7/httpclient"
	"net/http"
)

type HeimdallAgentClient struct {
	client *httpclient.Client
}

func (c *HeimdallAgentClient) Post(url string, contentType string, body []byte, compress bool) (*http.Response, error) {
	headers := http.Header{}
	headers.Set("Content-Type", contentType)
	if !compress {
		return c.client.Post(url, bytes.NewReader(body), headers)
	}

	cBody, err := Compress(body)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(cBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}

func NewHeimdallAgentClient() Client {
	return &HeimdallAgentClient{httpclient.NewClient()}
}
