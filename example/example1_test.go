package example

import (
	"fmt"
	"github.com/StudyPlace-io/work-pool-framework/pkg/workerpool"
	"k8s.io/klog/v2"
	"testing"
	"time"
)

/*
 使用方法：
 1. 创建工作池
 2. 定义需要的任务func
 3. 遍历任务数，放入全局队列
 4. 启动工作池
*/

func TestTaskPool1(t *testing.T) {

	// 建立一个工作池
	// input:池数量
	pool := workerpool.NewPool(5, workerpool.WithTimeout(1), workerpool.WithResultCallback(func(i interface{}) {
		fmt.Println("result: ", i)
	}))

	// 需要处理的任务
	tt := func(data interface{}) (interface{}, error) {
		taskID := data.(int)
		// 业务逻辑

		time.Sleep(100 * time.Millisecond)
		klog.Info("Task ", taskID, " processed")
		return nil, nil
	}

	// 准备多个个任务
	for i := 1; i <= 1000; i++ {

		// 需要做的任务
		task := workerpool.NewTaskInstance(fmt.Sprintf("task-%v", i), i, tt)

		// 所有的任务放入全局队列中
		pool.AddGlobalQueue(task)
	}
	pool.Run() // 启动

}
