package example

import (
	"fmt"
	"github.com/myconcurrencytools/workpoolframework/pkg/workerpool"
	"k8s.io/klog/v2"
	"math/rand"
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

func TestTaskPool2(t *testing.T) {

	// 建立一个池，

	// pool := workerpool.NewPool(5)
	pool := workerpool.NewPool(5, workerpool.WithMaxWorkerNum(25), workerpool.WithTimeout(1), workerpool.WithErrorCallback(func(err error) {
		if err != nil {
			fmt.Println("error handler: ", err)
		}
	}), workerpool.WithResultCallback(func(i interface{}) {
		fmt.Println("result: ", i)
	}))

	// 准备100个任务
	for i := 1; i <= 100; i++ {

		// 需要做的任务
		task := workerpool.NewTaskInstance(fmt.Sprintf("task-%v", i), i, func(data interface{}) (interface{}, error) {
			taskID := data.(int)

			/*
			   业务逻辑
			*/
			time.Sleep(100 * time.Millisecond)
			klog.Info("Task ", taskID, " processed")
			return nil, nil
		})

		// 所有的任务放入list中
		pool.AddGlobalQueue(task)
	}

	// 启动在后台等待执行
	go pool.RunBackground()

	for {
		taskID := rand.Intn(100) + 20

		//// 模拟一个退出条件
		if taskID%7 == 0 {
			klog.Info("taskID: ", taskID, "pool stop!")
			pool.StopBackground()
			break
		}

		time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
		// 模拟后续加入pool
		task := workerpool.NewTaskInstance(fmt.Sprintf("task-%v", taskID), taskID, func(data interface{}) (interface{}, error) {
			taskID := data.(int)
			time.Sleep(3 * time.Second)
			klog.Info("Task ", taskID, " processed")
			return nil, nil
		})

		pool.AddTask(task)
	}

	fmt.Println("finished...")
}
