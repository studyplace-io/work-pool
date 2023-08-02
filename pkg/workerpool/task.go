package workerpool

/*
 本质：用全局的切片分配任务给多个workers并发处理。
*/

// Task 任务接口，由工作池抽象出的具体执行单元，
// 当pool启动时，会从chan中不断读取Task接口对象执行
type Task interface {
	Execute() error
	GetTaskName() string
}

// TaskInstance 一个具体任务需求
type TaskInstance struct {
	Name string
	Err  error                   // 返回错误
	Data interface{}             // 真正的处理数据
	f    func(interface{}) error // 处理函数
}

// NewTaskInstance 建立任务
func NewTaskInstance(name string, data interface{}, f func(interface{}) error) *TaskInstance {
	return &TaskInstance{
		Name: name,
		Data: data,
		f:    f,
	}
}

func (t *TaskInstance) Execute() error {
	t.Err = t.f(t.Data) // 执行任务。如果任务执行错误，赋值err
	return nil
}

func (t *TaskInstance) GetTaskName() string {
	return t.Name
}

