package scheduler

import (
	"fmt"
	"github.com/myconcurrencytools/workpoolframework/pkg/workerpool"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	s := NewScheduler(5)

	s.Start()

	tsk := workerpool.NewTaskInstance("task1", "aaa", func(i interface{}) (interface{}, error) {
		fmt.Println(i)
		return nil, nil
	})

	s.AddTask(tsk)

	<-time.After(time.Second * 60)
	s.Stop()
}
