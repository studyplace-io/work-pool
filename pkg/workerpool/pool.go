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
func NewPool(concurrency int) *Pool {
	return &Pool{
		Tasks:         make([]*Task, 0),
		concurrency:   concurrency,
		collector:     make(chan *Task, 10),
		runBackground: make(chan bool),
	}
}

// AddGlobalQueue 加入工作池的全局队列，静态加入，用于启动工作池前的任务加入时使用，
// 在工作池启动后，推荐使用AddTask() 方法动态加入工作池
func (p *Pool) AddGlobalQueue(task *Task) {
	p.Tasks = append(p.Tasks, task)
}

func (p *Pool) Run() {
	// 总共会开启p.concurrency个goroutine （因为Start函数）
	for i := 1; i <= p.concurrency; i++ {
		worker := newWorker(p.collector, i)
		worker.start(&p.wg)
	}

	for len(p.Tasks) == 0 {
		klog.Error("no task in global queue...")
		time.Sleep(time.Millisecond)
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

// AddTask 把任务放入chan，当工作池启动后，动态加入使用
func (p *Pool) AddTask(task *Task) {
	// 放入chan
	p.collector <- task
}

// RunBackground 后台运行，需要启动一个goroutine来执行
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

	for len(p.Tasks) == 0 {
		klog.Error("no task in global queue...")
		time.Sleep(time.Millisecond)
	}

	for i := range p.Tasks {
		p.collector <- p.Tasks[i]
	}

	// 阻塞，等待关闭通知
	<-p.runBackground

}

// StopBackground 停止后台运行，需要chan通知
func (p *Pool) StopBackground() {
	klog.Info("pool close!")
	for i := range p.Workers {
		p.Workers[i].stop()
	}
	p.runBackground <- true
}
