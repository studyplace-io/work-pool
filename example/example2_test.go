package example

import (
	"github.com/myconcurrencytools/workpoolframework/pkg/workerpool"
	"k8s.io/klog/v2"
	"math/rand"
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

func TestTaskPool2(t *testing.T) {

	// 准备存放任务的地方
	var allTask []*workerpool.Task
	// 准备100个任务
	for i := 1; i <= 100; i++ {

		// 需要做的任务
		task := workerpool.NewTask(func(data interface{}) error {
			taskID := data.(int)

			/*
				业务逻辑
			*/
			time.Sleep(100 * time.Millisecond)
			klog.Info("Task ", taskID, " processed")
			return nil
		}, i)

		// 所有的任务放入list中
		allTask = append(allTask, task)
	}

	// 建立一个池，
	// input:待处理的任务对列;池数量

	pool := workerpool.NewPool(allTask, 5)

	//
	go func() {
		for {
			taskID := rand.Intn(100) + 20

			// 随意使用一个用例让pool停止
			if taskID%7 == 0 {
				klog.Info("taskID: ", taskID, "pool stop!")
				pool.Stop()
			}

			time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
			task := workerpool.NewTask(func(data interface{}) error {
				taskID := data.(int)
				time.Sleep(100 * time.Millisecond)
				klog.Info("Task ", taskID, " processed")
				return nil
			}, taskID)

			pool.AddTask(task)
		}

	}()

	// 确定所有任务后，才能调用
	pool.RunBackground()

}
