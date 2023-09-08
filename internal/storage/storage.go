package storage

import (
	"github.com/superles/yapmetrics/internal/storage/repositories"
)

var repo RepositoryInterface

func init() {
	repo = &repositories.MemoryRepository{}
}

func Store() RepositoryInterface {
	return repo
}
