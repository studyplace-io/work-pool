package workerpool

import (
	"k8s.io/klog/v2"
)

/*
 本质：用全局的切片分配任务给多个workers并发处理。
*/

// Task 一个具体任务需求
type Task struct {
	Err  error                   // 返回错误
	Data interface{}             // 真正的处理数据
	f    func(interface{}) error // 处理函数
}

// NewTask 建立任务
func NewTask(f func(interface{}) error, data interface{}) *Task {
	return &Task{
		Data: data,
		f:    f,
	}
}

// process 执行任务的函数。
func (t *Task) process(workerID int) {
	klog.Info("worker: ", workerID, ", processes task: ", t.Data)
	t.Err = t.f(t.Data) // 执行任务。如果任务执行错误，赋值err
}
