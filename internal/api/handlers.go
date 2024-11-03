package api

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	_ "workerPool/docs"
	"workerPool/internal/worker"
)

// PingHandler отвечает на запрос PING -> PONG
// @Summary Пинг-сервер
// @Description Проверка доступности сервера
// @Tags Health
// @Success 200 {string} string "PONG"
// @Router /api/ping [get]
func PingHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("PONG")); err != nil {
		log.Fatalf("failed to write response: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// AddJobHandler обрабатывает добавление задачи в пул
// @Summary Добавить задачу
// @Description Добавляет задачу в пул для обработки воркерами
// @Tags Worker
// @Param job query string true "Задача для обработки"
// @Success 200 {string} string "Задача добавлена"
// @Failure 400 {string} string "job is required"
// @Router /api/add_job [post]
func AddJobHandler(pool *worker.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		job := r.URL.Query().Get("job")
		if job == "" {
			http.Error(w, "job is required", http.StatusBadRequest)
			return
		}

		// Назначаем задание и получаем ID воркера
		workerID := pool.AssignJob(job)
		if workerID == -1 {
			http.Error(w, "Нет доступных воркеров для обработки задачи", http.StatusServiceUnavailable)
			return
		}

		// Пишем ответ с ID воркера
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprintf(w, "Задача добавлена: %s, ей занимается воркер с ID: %d\n", job, workerID)
		if err != nil {
			log.Printf("failed to write response: %v", err)
		}
	}
}

// AddWorkerHandler обрабатывает добавление нового воркера
// @Summary Добавить воркера
// @Description Добавляет нового воркера в пул
// @Tags Worker
// @Success 200 {string} string "Воркер добавлен с уникальным ID"
// @Router /api/add_worker [post]
func AddWorkerHandler(pool *worker.Pool, poolLock *sync.Mutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		poolLock.Lock()
		defer poolLock.Unlock()
		newID := len(pool.Workers) + 1
		pool.AddWorker(newID)
		_, err := fmt.Fprintf(w, "Воркер добавлен с ID: %d\n", newID)

		if err != nil {
			log.Fatalf("failed to write response: %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// RemoveWorkerHandler обрабатывает удаление воркера
// @Summary Удалить воркера
// @Description Удаляет одного воркера из пула
// @Tags Worker
// @Success 200 {string} string "Воркер удален"
// @Router /api/remove_worker [delete]
func RemoveWorkerHandler(pool *worker.Pool, poolLock *sync.Mutex) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		poolLock.Lock()
		defer poolLock.Unlock()
		if len(pool.Workers) == 0 {
			http.Error(w, "Нет воркеров для удаления", http.StatusBadRequest)
			return
		}
		pool.RemoveWorker()
		remainingWorkers := len(pool.Workers)
		_, err := fmt.Fprintln(w, "Воркер удален, количество живых воркеров: ", remainingWorkers)
		if err != nil {
			log.Fatalf("failed to write response: %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
