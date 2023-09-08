package storage

import (
	"github.com/superles/yapmetrics/internal/storage/repositories"
)

var repo RepositoryInterface

func init() {
	repo = &repositories.MemoryRepository{}
	var slice []any
	slice = append(slice, 1)
	slice = append(slice, 1.11)
}

func Store() RepositoryInterface {
	return repo
}
