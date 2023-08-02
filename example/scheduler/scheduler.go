package scheduler

import (
	"github.com/myconcurrencytools/workpoolframework/pkg/workerpool"
	"log"
)

type Scheduler struct {
	pool *workerpool.Pool
}

func NewScheduler(workerNum int) *Scheduler {
	s := &Scheduler{
		pool: workerpool.NewPool(workerNum),
	}
	return s
}

func (s *Scheduler) Start() {
	go func() {
		log.Printf("start scheduler...")
		s.pool.RunBackground()
	}()
}

func (s *Scheduler) Stop() {
	s.pool.StopBackground()
}

func (s *Scheduler) AddTask(task workerpool.Task) {
	s.pool.AddTask(task)
}
