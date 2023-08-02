package workerpool

import (
	"k8s.io/klog/v2"
	"math/rand"
	"sync"
	"time"
)

// Pool 工作池
type Pool struct {
	// list 装task
	Tasks   []Task
	// Workers 列表
	Workers []*worker
	// 工作池数量
	concurrency int
	// collector 用来输入所有Task对象的chan
	collector chan Task
	// runBackground 后台运行时，结束时需要传入的标示
	runBackground  chan bool
	// timeout 超时时间
	timeout time.Duration
	// errorCallback 当任务发生错误时的回调方法
	errorCallback func(err error)
	// resultCallback 当任务有结果时的回调方法
	resultCallback func(result interface{})
	wg             sync.WaitGroup
}

// NewPool 建立一个pool
func NewPool(concurrency int, opts ...Option) *Pool {
	p := &Pool{
		Tasks:         make([]Task, 0),
		Workers:       make([]*worker, 0),
		concurrency:   concurrency,
		collector:     make(chan Task, 10),
		runBackground: make(chan bool),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// AddGlobalQueue 加入工作池的全局队列，静态加入，用于启动工作池前的任务加入时使用，
// 在工作池启动后，推荐使用AddTask() 方法动态加入工作池
func (p *Pool) AddGlobalQueue(task Task) {
	p.Tasks = append(p.Tasks, task)
}

// Run 启动pool，使用Run()方法调用时，只能使用AddGlobalQueue加入全局队列，
// 一旦Run启动后，就不允许调用AddTask加入Task，如果需动态加入pool，可以使用
// RunBackground方法
func (p *Pool) Run() {
	// 总共会开启p.concurrency个goroutine
	// 启动pool中的每个worker都传入collector chan
	for i := 1; i <= p.concurrency; i++ {
		wr := newWorker(i, p.timeout, p.errorCallback, p.resultCallback)
		p.Workers = append(p.Workers, wr)
		wr.start(&p.wg)
	}

	// 如果全局队列没任务，提示一下
	if len(p.Tasks) == 0 {
		klog.Info("no task in global queue...")
	}

	go p.dispatch()

	// 把放在tasks列表的的任务放入collector
	for i := range p.Tasks {
		p.collector <- p.Tasks[i]
	}

	// 注意，这里需要close chan。
	close(p.collector)

	// 阻塞，等待所有的goroutine执行完毕
	p.wg.Wait()
}

// dispatch 由pool chan中不断分发给worker chan
// 使用随机分配的方式
func (p *Pool) dispatch() {
	for task := range p.collector {
		index := rand.Intn(p.concurrency)
		p.Workers[index].taskChan <- task
	}

	for _, v := range p.Workers {
		close(v.taskChan)
	}

}

// AddTask 把任务放入chan，当工作池启动后，动态加入使用
func (p *Pool) AddTask(task Task) {
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
		wk := newWorker(i, p.timeout, p.errorCallback, p.resultCallback)
		p.Workers = append(p.Workers, wk)

		go wk.startBackground()
	}

	go p.dispatch()

	// 如果全局队列没任务，提示一下
	if len(p.Tasks) == 0 {
		klog.Info("no task in global queue...")
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
