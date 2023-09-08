package storage

import "github.com/superles/yapmetrics/internal/types"

type RepositoryInterface interface {
	Get() types.Metric
	Add(types.Metric)
}
