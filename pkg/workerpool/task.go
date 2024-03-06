package workerpool

import "fmt"

/*
 本质：用全局的切片分配任务给多个workers并发处理。
*/

// Task 任务接口，由工作池抽象出的具体执行单元，
// 当 workpool 启动时，会从 chan 中不断读取 Task接口对象 执行
type Task interface {
	// Execute 执行任务方法
	Execute() (interface{}, error)
	// GetTaskName 获取任务名
	GetTaskName() string
}

// TaskInstance 一个具体任务需求
type TaskInstance struct {
	Name string
	// Err 返回错误
	Err error
	// Data 真正的处理数据
	Data interface{}
	// f 处理任务函数，由调用方传入
	f TaskFunc
}

type TaskFunc func(interface{}) (interface{}, error)

// NewTaskInstance 建立任务
func NewTaskInstance(name string, data interface{}, f TaskFunc) *TaskInstance {
	return &TaskInstance{
		Name: name,
		Data: data,
		f:    f,
	}
}

func (t *TaskInstance) Execute() (interface{}, error) {
	// 判断是否有传入 func
	if t.f == nil {
		return nil, fmt.Errorf("no task func init")
	}
	// 执行任务。如果任务执行错误，赋值err
	result, err := t.f(t.Data)
	if err != nil {
		t.Err = err
		return nil, err
	}
	return result, nil
}

func (t *TaskInstance) GetTaskName() string {
	return t.Name
}
