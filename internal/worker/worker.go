package worker

import (
	"fmt"
	"sync"
)

type Worker struct {
	ID       int
	JobQueue chan string
	Quit     chan bool
}

func (w *Worker) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case job := <-w.JobQueue:
			fmt.Printf("Воркер[%d] обрабатывает данные: %s\n", w.ID, job)
		case <-w.Quit:
			fmt.Printf("Воркер[%d] завершил свою работу.\n", w.ID)
			return
		}
	}
}
