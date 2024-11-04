package test

import (
	"testing"
	"time"
	"workerPool/internal/worker"
)

func TestAddWorker(t *testing.T) {
	pool := worker.NewPool()
	pool.AddWorker(1)

	if len(pool.Workers) != 1 {
		t.Errorf("Expecting 1 worker, got %d", len(pool.Workers))
	}

	if pool.Workers[0].ID != 1 {
		t.Errorf("Expecting worker ID to be 1, got %d", pool.Workers[0].ID)
	}
}

func TestRemoveWorker(t *testing.T) {
	pool := worker.NewPool()
	pool.AddWorker(1)
	pool.AddWorker(2)

	pool.RemoveWorker()
	if len(pool.Workers) != 1 {
		t.Errorf("Expecting 1 worker, got %d", len(pool.Workers))
	}
	if pool.Workers[0].ID != 1 {
		t.Errorf("Expecting worker ID to be 1, got %d", pool.Workers[0].ID)
	}
}

func TestWait(t *testing.T) {
	pool := worker.NewPool()
	pool.AddWorker(1)
	pool.AddWorker(2)

	go func() {
		time.Sleep(100 * time.Millisecond)
		pool.RemoveWorker()
		pool.RemoveWorker()
	}()

	done := make(chan struct{})
	go func() {
		pool.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Errorf("Expecting Wait to finish after all workers are removed, but it timed out")
	}
}
