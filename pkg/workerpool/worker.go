package workerpool

import (
	"context"
	"fmt"
	"k8s.io/klog/v2"
	"sync"
	"time"
)

// worker 执行任务的消费者
type worker struct {
	ID int // 消费者的id
	// taskChan 等待处理的任务chan (每个worker都有一个自己的chan)
	taskChan chan Task
	// quit 停止通知
	quit chan bool
	// timeout 超时时间
	timeout time.Duration
	// errorCallback 当任务发生错误时的回调方法
	errorCallback func(err error)
	// resultCallback 当任务有结果时的回调方法
	resultCallback func(result interface{})
}

// newWorker 创建worker
func newWorker(ID int, timeout time.Duration, errorCallback func(err error), resultCallback func(interface{})) *worker {
	return &worker{
		ID:             ID,
		taskChan:       make(chan Task, 10),
		quit:           make(chan bool),
		timeout:        timeout,
		errorCallback:  errorCallback,
		resultCallback: resultCallback,
	}
}

// executeTask 执行任务
func (wr *worker) executeTask(task Task) (interface{}, error) {
	var err error
	var result interface{}
	if wr.timeout > 0 {
		result, err = wr.executeTaskWithTimeout(task)
	} else {
		result, err = wr.executeTaskWithoutTimeout(task)
	}
	return result, err
}

// executeTaskWithTimeout 执行任务有超时的情况
func (wr *worker) executeTaskWithTimeout(task Task) (interface{}, error) {

	ctx, cancel := context.WithTimeout(context.Background(), wr.timeout*time.Second)
	defer cancel()

	var result interface{}
	var err error
	done := make(chan struct{})

	// 异步执行并等待
	go func() {
		result, err = task.Execute()
		close(done)
	}()

	// 阻塞等待超时先到还是任务先执行完成
	select {
	case <-done:
		return result, err
	case <-ctx.Done():
		return nil, fmt.Errorf("task timed out...")
	}
}

// executeTaskWithoutTimeout 执行任务
func (wr *worker) executeTaskWithoutTimeout(task Task) (interface{}, error) {
	return task.Execute()
}

// start 执行任务，遍历taskChan，每个worker都启一个goroutine执行。
func (wr *worker) start(wg *sync.WaitGroup) {
	klog.Infof("Starting worker: %v\n", wr.ID)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// 不断从chan中取出task执行
		for task := range wr.taskChan {
			if task == nil {
				continue
			}
			klog.Info("worker: ", wr.ID, ", processes task: ", task.GetTaskName())
			result, err := wr.executeTask(task)
			wr.handleResult(result, err)
		}
	}()
}

// startBackground 后台执行
func (wr *worker) startBackground() {
	klog.Info("Starting worker background: ", wr.ID)
	for {
		select {
		case task := <-wr.taskChan:
			if task == nil {
				continue
			}
			klog.Info("worker: ", wr.ID, ", processes task: ", task.GetTaskName())
			result, err := wr.executeTask(task)
			wr.handleResult(result, err)
		case <-wr.quit:
			return
		}
	}
}

// handleResult 处理任务结束的方法
func (wr *worker) handleResult(result interface{}, err error) {
	if err != nil && wr.errorCallback != nil {
		wr.errorCallback(err)
	} else if wr.resultCallback != nil {
		wr.resultCallback(result)
	}
}

// stop 停止worker
func (wr *worker) stop() {
	klog.Info("Closing worker: ", wr.ID)
	go func() {
		wr.quit <- true
	}()
}