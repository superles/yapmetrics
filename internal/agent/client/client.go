package client

import (
	"context"
	"github.com/superles/yapmetrics/internal/metric"
)

type Client interface {
	Send(ctx context.Context, endpoint string, metrics []metric.Metric) error
}
