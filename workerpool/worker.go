package workerpool

import (
	"fmt"
	"sync"
)

//
type Worker struct {
	// 消费者的id
	ID			int
	// 等待处理的任务chan (每个worker都有一个自己的chan)
	taskChan	chan *Task
	// 停止通知
	quit 	chan bool
}

// 建立新的消费者
func NewWorker(channel chan *Task, ID int) *Worker {
	return &Worker{
		ID: ID,
		taskChan: channel,
		quit: make(chan bool),
	}
}


// 执行，遍历taskChan，每个任务都启一个goroutine执行。
func (wr *Worker) Start(wg *sync.WaitGroup) {
	fmt.Printf("Starting worker %d\n", wr.ID)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for task := range wr.taskChan {
			process(wr.ID, task)
		}
	}()
}


func (wr *Worker) StartBackground() {
	fmt.Printf("Starting worker %d\n", wr.ID)

	for {
		select {
		case task := <- wr.taskChan:
			process(wr.ID, task)
		case <- wr.quit:
			return
		}
	}

}

func (wr *Worker) Stop() {
	fmt.Printf("Closing worker %d\n", wr.ID)

	go func() {
		wr.quit <- true
	}()
}

