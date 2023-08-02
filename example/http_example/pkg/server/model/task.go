package model

import "fmt"

const (
	TaskRunning = "running"
	TaskFail    = "fail"
	TaskSuccess = "success"
)

type MyTask struct {
	TaskName string      `json:"taskName"`
	TaskType string      `json:"taskType"`
	Input    interface{} `json:"input"`
	f        func(data interface{}) error
	Err      error
	Status   string
}

func (my *MyTask) ChooseTaskType() {
	if my.TaskType == "string" {
		my.f = func(data interface{}) error {
			fmt.Println("string task...", data)
			return nil
		}
	} else {
		my.f = func(data interface{}) error {
			fmt.Println("int task...", data)
			return nil
		}
	}
}

func (my *MyTask) Execute() error {

	my.Status = TaskRunning

	if err := my.f(my.Input); err != nil {
		my.Err = err
		my.Status = TaskFail
		return err
	}
	my.Status = TaskSuccess
	return nil
}

func (my *MyTask) GetTaskName() string {
	return my.TaskName
}
