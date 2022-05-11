package workerpool

import "fmt"


/*
	本质：用全局的切片分配任务给多个workers并发处理。



 */



// 一个具体任务需求
type Task struct {
	// 返回错误
	Err		error
	// 真正的处理数据
	Data 	interface{}
	// 处理函数
	f 		func(interface{}) error
}

// 建立任务
func NewTask(f func(interface{}) error, data interface{}) *Task {
	return &Task{
		Data: data,
		f: f,
	}
}

// 执行任务的函数。
func process(workerID int, task *Task) {
	fmt.Printf("worker %d processes task %v\n", workerID, task.Data)
	task.Err = task.f(task.Data)	// 执行
}

