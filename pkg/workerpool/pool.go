package workerpool

import (
	"k8s.io/klog/v2"
	"sync"
	"time"
)

// Pool 工作池
type Pool struct {
	// list 装task
	Tasks   []*Task
	Workers []*worker

	// 工作池数量
	concurrency int
	// 用来装
	collector     chan *Task
	runBackground chan bool
	wg            sync.WaitGroup
}

// NewPool 建立一个pool
func NewPool(tasks []*Task, concurrency int) *Pool {
	return &Pool{
		Tasks:         tasks,
		concurrency:   concurrency,
		collector:     make(chan *Task, 10),
		runBackground: make(chan bool),
	}
}

func (p *Pool) Run() {
	// 总共会开启p.concurrency个goroutine （因为Start函数）
	for i := 1; i <= p.concurrency; i++ {
		worker := newWorker(p.collector, i)
		worker.start(&p.wg)
	}

	// 把好的任务放入collector
	for i := range p.Tasks {
		p.collector <- p.Tasks[i]

	}

	// 关闭通道
	close(p.collector)

	// 阻塞，等待所有的goroutine执行完毕
	p.wg.Wait()
}

// 把任务放入chan
func (p *Pool) AddTask(task *Task) {
	// 放入chan
	p.collector <- task
}

func (p *Pool) RunBackground() {
	// 启动goroutine，打印。
	go func() {
		for {
			klog.Info("Waiting for tasks to come in... \n")
			time.Sleep(10 * time.Second)
		}
	}()

	// 启动workers 数量： p.concurrency
	for i := 1; i <= p.concurrency; i++ {
		workers := newWorker(p.collector, i)
		p.Workers = append(p.Workers, workers)

		go workers.startBackground()
	}

	for i := range p.Tasks {
		p.collector <- p.Tasks[i]
	}

	// 阻塞，等待关闭通知
	<-p.runBackground

}

func (p *Pool) StopBackground() {
	klog.Info("pool close!")
	for i := range p.Workers {
		p.Workers[i].stop()
	}
	p.runBackground <- true
}
