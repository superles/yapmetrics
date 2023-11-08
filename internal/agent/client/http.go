package client

import (
	"bytes"
	"github.com/superles/yapmetrics/internal/utils/hasher"
	"net/http"
)

type AgentClientParams struct {
	Key string
}

type HTTPAgentClient struct {
	client *http.Client
	params AgentClientParams
}

func (c *HTTPAgentClient) Post(url string, contentType string, body []byte, compress bool) (*http.Response, error) {
	var hash string
	if len(c.params.Key) != 0 {
		hash = hasher.Encode(body, []byte(c.params.Key))
	}
	if compress {
		var err error
		body, err = Compress(body)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	if compress {
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")
	}
	if len(c.params.Key) != 0 {
		req.Header.Set("HashSHA256", hash)
	}
	return c.client.Do(req)
}

func NewHTTPAgentClient(params AgentClientParams) Client {
	return &HTTPAgentClient{&http.Client{}, params}
}
