package model

type MockMap map[string]*MyTask

var TaskMap MockMap

func init() {
	TaskMap = MockMap{}
}
