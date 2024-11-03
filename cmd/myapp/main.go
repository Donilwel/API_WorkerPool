package main

import (
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"sync"
	_ "workerPool/docs" // Подключаем документацию Swagger только в main.go
	"workerPool/internal/api"
	"workerPool/internal/config"
	"workerPool/internal/worker"
)

var (
	pool     *worker.Pool
	poolLock sync.Mutex
)

// @title Worker Pool API
// @version 1.0
// @description Пример API для управления worker-пулом.
// @host localhost:8080
// @BasePath /api
func main() {
	// Загружаем переменные окружения
	config.LoadEnv()

	// Инициализируем пул воркеров
	pool = worker.NewPool()
	for i := 1; i <= 3; i++ {
		pool.AddWorker(i)
	}

	// Создаем маршрутизатор
	r := mux.NewRouter()
	apiRouter := r.PathPrefix("/api").Subrouter()

	// Настраиваем обработчики API
	apiRouter.HandleFunc("/ping", api.PingHandler).Methods("GET")
	apiRouter.HandleFunc("/add_worker", api.AddWorkerHandler(pool, &poolLock)).Methods("POST")
	apiRouter.HandleFunc("/add_job", api.AddJobHandler(pool)).Methods("POST")
	apiRouter.HandleFunc("/remove_worker", api.RemoveWorkerHandler(pool, &poolLock)).Methods("DELETE")

	// Добавляем маршрут для Swagger UI
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Запускаем сервер
	log.Printf("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

	// Ожидаем завершения всех горутин
	pool.Wait()
}
