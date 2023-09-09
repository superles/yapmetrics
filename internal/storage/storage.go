package storage

import "github.com/superles/yapmetrics/internal/storage/repository"

var MetricRepository = new(repository.MemoryMetricRepository)

func init() {
	MetricRepository = new(repository.MemoryMetricRepository)
}
