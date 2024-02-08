package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mailru/easyjson"
	"github.com/superles/yapmetrics/internal/metric"
	"github.com/superles/yapmetrics/internal/utils/encoder"
	"github.com/superles/yapmetrics/internal/utils/hasher"
	"net/http"
)

type AgentClientParams struct {
	Key      string
	RealIP   string
	Compress bool
	Encoder  *encoder.Encoder
}

type HTTPAgentClient struct {
	client *http.Client
	params AgentClientParams
}

func stringifyMetrics(metrics []metric.Metric) ([]byte, error) {
	var col metric.JSONDataCollection
	for _, item := range metrics {
		updatedJSON, err := item.ToJSON()
		if err != nil {
			return nil, err
		}
		col = append(col, updatedJSON)
	}
	rawBytes, err := easyjson.Marshal(col)

	if err != nil {
		return nil, err
	}
	return rawBytes, nil
}

func (c *HTTPAgentClient) Send(ctx context.Context, endpoint string, metrics []metric.Metric) error {

	url := fmt.Sprintf("http://%s/updates/", endpoint)
	contentType := "application/json"

	body, err := stringifyMetrics(metrics)

	if err != nil {
		return err
	}

	if c.params.Encoder != nil {
		if body, err = c.params.Encoder.Encrypt(body); err != nil {
			return err
		}
	}

	var hash string
	if len(c.params.Key) != 0 {
		hash = hasher.Encode(body, []byte(c.params.Key))
	}

	if c.params.Compress {
		body, err = Compress(body)
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	if c.params.Compress {
		req.Header.Set("Content-Encoding", "gzip")
		req.Header.Set("Accept-Encoding", "gzip")
	}
	if len(c.params.RealIP) > 0 {
		req.Header.Set("X-Real-Ip", c.params.RealIP)
	}
	if len(c.params.Key) != 0 {
		req.Header.Set("HashSHA256", hash)
	}

	response, err := c.client.Do(req)

	if err != nil {
		return err
	}

	return response.Body.Close()
}

func NewHTTPAgentClient(params AgentClientParams) Client {
	return &HTTPAgentClient{&http.Client{}, params}
}
