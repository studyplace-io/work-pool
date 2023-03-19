package example

import (
	"github.com/myconcurrencytools/workpoolframework/pkg/workerpool"
	"k8s.io/klog/v2"
	"testing"
	"time"
)

/*
	使用方法：
	1. 准备全局的任务队列，用于存放任务
	2. 定义需要的任务func
	3. 遍历任务数，放入全局队列
	4. 创建且启动工作池
 */

func TestTaskPool1(t *testing.T) {


	// 准备存放任务的地方，全局任务队列
	var allTask []*workerpool.Task

	// 需要处理的任务
	tt := func(data interface{}) error {
		taskID := data.(int)
		// 业务逻辑

		time.Sleep(100 * time.Millisecond)
		klog.Info("Task ", taskID, " processed")
		return nil
	}

	// 准备多个个任务
	for i := 1; i <= 1000; i++ {

		// 需要做的任务
		task := workerpool.NewTask(tt, i)

		// 所有的任务放入全局队列中
		allTask = append(allTask, task)
	}

	// 建立一个工作池
	// input:待处理的任务对列;池数量
	pool := workerpool.NewPool(allTask, 5)
	pool.Run() // 启动

}
