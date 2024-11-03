package worker

import (
	"fmt"
	"math/rand"
	"sync"
)

type Pool struct {
	Workers []*Worker
	wg      *sync.WaitGroup
}

func NewPool() *Pool {
	return &Pool{
		Workers: make([]*Worker, 0),
		wg:      &sync.WaitGroup{},
	}
}

func (p *Pool) AddWorker(id int) {
	worker := &Worker{
		ID:       id,
		JobQueue: make(chan string),
		Quit:     make(chan bool),
	}
	p.Workers = append(p.Workers, worker)
	p.wg.Add(1)
	go worker.Start(p.wg)
}

func (p *Pool) RemoveWorker() {
	if len(p.Workers) == 0 {
		fmt.Println("worker pool is empty")
		return
	}
	worker := p.Workers[len(p.Workers)-1]
	worker.Quit <- true
	p.Workers = p.Workers[:len(p.Workers)-1]
}

func (p *Pool) AssignJob(job string) {
	if len(p.Workers) == 0 {
		fmt.Println("Нет доступных воркеров для обработки задачи")
		return
	}
	randWorker := p.Workers[rand.Intn(len(p.Workers))]
	randWorker.JobQueue <- job
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
