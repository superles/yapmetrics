package main

import (
	"github.com/superles/yapmetrics/internal/agent"
	"github.com/superles/yapmetrics/internal/storage/memstorage"
)

func main() {
	storage := memstorage.New()
	agent.New(storage).Run()
}
