package client

import (
	"bytes"
	"github.com/superles/yapmetrics/internal/utils/hasher"
	"net/http"
)

type HTTPAgentClient struct {
	client *http.Client
	key    string
}

func (c *HTTPAgentClient) Post(url string, contentType string, body []byte, compress bool) (*http.Response, error) {
	var hash string
	if len(c.key) != 0 {
		hash = hasher.Encode(string(body), c.key)
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
	if len(c.key) != 0 {
		req.Header.Set("HashSHA256", hash)
	}
	return c.client.Do(req)
}

func NewHTTPAgentClient(key string) Client {
	return &HTTPAgentClient{&http.Client{}, key}
}
