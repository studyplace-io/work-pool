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
	Tasks []Task
	// Workers 列表
	Workers []*worker
	// concurrency 工作池数量
	concurrency int
	// maxWorkerNum 最大工作池数量
	maxWorkerNum int
	// collector 用来输入所有Task对象的chan
	collector chan Task
	// runBackground 后台运行时，结束时需要传入的标示
	runBackground chan bool
	// timeout 超时时间
	timeout time.Duration
	// errorCallback 当任务发生错误时的回调方法
	errorCallback func(err error)
	// resultCallback 当任务有结果时的回调方法
	resultCallback func(result interface{})
	wg             sync.WaitGroup
	lock           sync.Mutex
}

const (
	defaultMaxWorkerNum = 20
)

// NewPool 建立一个pool
func NewPool(concurrency int, opts ...Option) *Pool {
	p := &Pool{
		Tasks:         make([]Task, 0),
		Workers:       make([]*worker, 0),
		concurrency:   concurrency,
		collector:     make(chan Task, 500),
		runBackground: make(chan bool),
		lock:          sync.Mutex{},
	}

	if p.maxWorkerNum == 0 {
		p.maxWorkerNum = defaultMaxWorkerNum
	}

	if p.concurrency > p.maxWorkerNum {
		p.maxWorkerNum = p.concurrency
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

	// 增加扩展功能
	go p.dispatch()
	go p.autoScale()

	// 把放在tasks列表的的任务放入collector
	for i := range p.Tasks {
		p.collector <- p.Tasks[i]
	}

	// 注意，这里需要close chan。
	close(p.collector)

	// 阻塞，等待所有的goroutine执行完毕
	p.wg.Wait()
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

	// 增加扩展功能
	go p.dispatch()
	go p.autoScale()

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
	close(p.collector)
	for _, k := range p.Workers {
		k.stop()
	}
	//p.runBackground <- true
}

// dispatch 由pool chan中不断分发给worker chan
// 使用随机分配的方式
func (p *Pool) dispatch() {
	for task := range p.collector {
		p.lock.Lock()
		index := rand.Intn(len(p.Workers))
		p.Workers[index].taskChan <- task
		p.lock.Unlock()
	}
	// 当p.collector被关闭时，代表任务都执行完毕
	// 需要把所有worker的chan都关闭
	for _, v := range p.Workers {
		close(v.taskChan)
	}

}

// autoScale 监测自动扩缩容
func (p *Pool) autoScale() {
	for {

		time.Sleep(5 * time.Second)

		p.lock.Lock()

		numWorkers := len(p.Workers)
		numJobs := len(p.collector)

		// 如果全局chan中任务数为0，把数量设置为最小worker数
		if numJobs == 0 && len(p.Workers) > p.concurrency {
			p.scaleWorkers(p.concurrency)
		} else {
			if numWorkers == 0 {
				continue
			}
			// 如果全局队列内的任务数量大于总容量的3/4，
			// 就认为任务堆积，需要扩容worker
			if numJobs > cap(p.collector)*3/4 {
				p.scaleWorkers(numWorkers + 1)
			}

		}

		p.lock.Unlock()
	}
}

// scaleWorkers 调整worker数方法方法
func (p *Pool) scaleWorkers(numWorkers int) {
	// 获取目前的worker数量
	currentNumWorkers := len(p.Workers)

	// 如果数量相等，直接return
	if currentNumWorkers == numWorkers {
		return
	}

	// 如果计算出的数量超过最大数量，直接使用最大数量
	if numWorkers >= p.maxWorkerNum {
		numWorkers = p.maxWorkerNum
	}

	// 如果期望数量比目前数量大，代表需要扩容
	if currentNumWorkers <= numWorkers {
		diff := numWorkers - currentNumWorkers
		for i := 0; i < diff; i++ {
			nwk := newWorker(len(p.Workers)+1, p.timeout, p.errorCallback, p.resultCallback)
			p.Workers = append(p.Workers, nwk)
			nwk.start(&p.wg)
		}
		klog.Infof("Scaled up %d workers\n", diff)
		klog.Infof("Scaled up to %d workers\n", numWorkers)
	} else {

		diff := currentNumWorkers - numWorkers
		for i := diff; i >= 0; i-- {
			nwk := p.Workers[i]
			nwk.stop()
			p.Workers = p.Workers[:i]
		}
		klog.Infof("Scaled down %d workers\n", diff)
		klog.Infof("Scaled down to %d workers\n", numWorkers)

	}
}