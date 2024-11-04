package test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"workerPool/internal/api"
	"workerPool/internal/worker"
)

func TestPingHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/ping", nil)
	w := httptest.NewRecorder()

	api.PingHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("got %d, want %d", resp.StatusCode, http.StatusOK)
	}
	body := w.Body.String()
	if body != "PONG" {
		t.Errorf("got %s, want PONG", body)
	}
}

func TestAddJobHandler(t *testing.T) {
	pool := worker.NewPool()
	pool.AddWorker(1)
	handler := api.AddJobHandler(pool)

	t.Run("TestWithoutParameters", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/add_job?job=", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
		}
	})

	t.Run("TestWithParameters", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/add_job?job=SomethingJob", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		body := rr.Body.String()
		example := "Задача добавлена: SomethingJob, ей занимается воркер с ID: 1\n"
		if body != example {
			t.Errorf("handler returned unexpected body: got\n %v want\n %v", body, example)
		}
	})
}

func TestAddWorkerHandler(t *testing.T) {
	pool := worker.NewPool()
	mu := sync.Mutex{}

	handler := api.AddWorkerHandler(pool, &mu)

	t.Run("AddWorker", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/add_worker", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		example := "Воркер добавлен с ID: 1\n"
		body := rr.Body.String()

		if body != example {
			t.Errorf("handler returned unexpected body: got\n %v want\n %v", body, example)
		}

		if len(pool.Workers) != 1 {
			t.Errorf("handler returned wrong len of workers: got %d, want %d", len(pool.Workers), 1)
		}
	})

	t.Run("AddSecondWorker", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/add_worker", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		example := "Воркер добавлен с ID: 2\n"
		body := rr.Body.String()

		if body != example {
			t.Errorf("handler returned unexpected body: got\n %v want\n %v", body, example)
		}

		if len(pool.Workers) != 2 {
			t.Errorf("handler returned wrong len of workers: got %d, want %d", len(pool.Workers), 2)
		}
	})
}

func TestRemoveWorkerHandler(t *testing.T) {
	pool := worker.NewPool()
	mu := sync.Mutex{}
	handler := api.RemoveWorkerHandler(pool, &mu)

	t.Run("RemoveWorkerWithZero", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/remove_worker", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		example := "Нет воркеров для удаления\n"
		body := rr.Body.String()

		if body != example {
			t.Errorf("handler returned unexpected body: got\n %v want\n %v", body, example)
		}

		if len(pool.Workers) != 0 {
			t.Errorf("handler returned wrong len of workers: got %d, want %d", len(pool.Workers), 0)
		}
	})
	pool.AddWorker(1)

	t.Run("RemoveWorkerWithOne", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/remove_worker", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		example := "Воркер удален, количество живых воркеров:  0\n"
		body := rr.Body.String()
		if body != example {
			t.Errorf("handler returned unexpected body: got\n %v want\n %v", body, example)
		}
		if len(pool.Workers) != 0 {
			t.Errorf("handler returned wrong len of workers: got %d, want %d", len(pool.Workers), 0)
		}
	})
	pool.AddWorker(2)

	t.Run("RemoveWorkerWithTwo", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/api/remove_worker", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		example := "Воркер удален, количество живых воркеров:  1\n"
		body := rr.Body.String()
		if body != example {
			t.Errorf("handler returned unexpected body: got\n %v want\n %v", body, example)
		}
		if len(pool.Workers) != 1 {
			t.Errorf("handler returned wrong len of workers: got %d, want %d", len(pool.Workers), 1)
		}
	})
}
