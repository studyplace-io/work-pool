package workerpool

import (
	"fmt"
	"sync"
	"time"
)

type Pool struct {
	// list 装task
	Tasks	[]*Task
	Workers []*Worker

	// 池的数量
	concurrency	int
	// 用来装
	collector	chan *Task
	runBackground chan bool
	wg 			sync.WaitGroup
}

// 建立一个pool
func NewPool(tasks []*Task, concurrency int) *Pool {
	return &Pool{
		Tasks: tasks,
		concurrency: concurrency,
		collector: make(chan *Task, 10),
		runBackground: make(chan bool),
	}
}

func (p *Pool) Run() {
	// 总共会开启p.concurrency个goroutine （因为Start函数）
	for i := 1; i <= p.concurrency; i++ {
		worker := NewWorker(p.collector, i)
		worker.Start(&p.wg)
	}

	// 把好的任务放入collector
	for i := range p.Tasks {
		p.collector <- p.Tasks[i]

	}

	// 关闭通道
	close(p.collector)

	p.wg.Wait()
}

// 把任务放入chan
func (p *Pool) AddTask(task *Task) {
	// 放入chan
	p.collector <- task
}


func (p *Pool) RunBackground() {
	go func() {
		for {
			fmt.Printf("Waiting for tasks to come in... \n")
			time.Sleep(10 * time.Second)
		}
	}()

	// 启动workers 数量： p.concurrency
	for i := 1; i < p.concurrency; i++ {
		workers := NewWorker(p.collector, p.concurrency)
		p.Workers = append(p.Workers, workers)

		go workers.StartBackground()
	}

	for i := range p.Tasks {
		p.collector <- p.Tasks[i]
	}

	<- p.runBackground

}


func (p *Pool) Stop() {

	for i := range p.Workers {
		p.Workers[i].Stop()
	}
	p.runBackground <- true
}


