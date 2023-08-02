package workerpool

import (
	"k8s.io/klog/v2"
	"sync"
)

// worker 执行任务的消费者
type worker struct {
	ID int // 消费者的id
	// 等待处理的任务chan (每个worker都有一个自己的chan)
	taskChan chan Task
	// 停止通知
	quit chan bool
}

// newWorker 建立新的消费者
func newWorker(channel chan Task, ID int) *worker {
	return &worker{
		ID:       ID,
		taskChan: channel,
		quit:     make(chan bool),
	}
}

// start 执行任务，遍历taskChan，每个worker都启一个goroutine执行。
func (wr *worker) start(wg *sync.WaitGroup) {
	klog.Info("Starting worker: ", wr.ID)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// 不断从chan中取出task执行
		for task := range wr.taskChan {
			klog.Info("worker: ", wr.ID, ", processes task: ", task.GetTaskName())
			task.Execute()
		}
	}()
}

// startBackground 后台执行
func (wr *worker) startBackground() {
	klog.Info("Starting worker background: ", wr.ID)
	for {
		select {
		case task := <-wr.taskChan:
			klog.Info("worker: ", wr.ID, ", processes task: ", task.GetTaskName())
			task.Execute()
		case <-wr.quit:
			return
		}
	}

}

func (wr *worker) stop() {
	klog.Info("Closing worker: ", wr.ID)
	go func() {
		wr.quit <- true
	}()
}