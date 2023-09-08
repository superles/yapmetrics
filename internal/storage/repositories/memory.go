package repositories

import (
	"github.com/superles/yapmetrics/internal/types"
)

type MemoryRepository struct {
	store []types.Metric
}

func (r *MemoryRepository) Add(item types.Metric) {
	r.store = append(r.store, item)
}

func (r *MemoryRepository) Get() types.Metric {
	if r.store == nil || len(r.store) == 0 {
		return types.Metric{}
	}
	return r.store[len(r.store)-1]
}
