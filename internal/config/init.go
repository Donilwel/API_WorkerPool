package config

import (
	"math/rand"
	"time"
	"workerPool/internal/worker"
)

func init() {
	pool := worker.NewPool()
	for i := 1; i <= 3; i++ {
		pool.AddWorker(i)
	}
	rand.Seed(time.Now().UnixNano())
}
