package client

import (
	"bytes"
	"github.com/gojek/heimdall/v7/httpclient"
	"net/http"
)

type HeimdallAgentClient struct {
	client *httpclient.Client
}

func (c *HeimdallAgentClient) Post(url string, contentType string, body []byte) (*http.Response, error) {
	headers := http.Header{}
	headers.Set("Content-Type", contentType)
	return c.client.Post(url, bytes.NewReader(body), headers)
}

func NewHeimdallAgentClient() Client {
	return &HeimdallAgentClient{httpclient.NewClient()}
}
