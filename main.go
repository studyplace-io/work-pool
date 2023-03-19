package main


import (
	"fmt"
	"github.com/myconcurrencytools/workpoolframework/workerpool"
	"k8s.io/klog/v2"
	"math/rand"
	"time"
)

func main() {

	// 准备存放任务的地方
	var allTask []*workerpool.Task
	// 准备100个任务
	for i := 1; i <= 71; i++ {

		// 需要做的任务
		task := workerpool.NewTask(func(data interface{}) error {
			taskID := data.(int)

			/*
				业务逻辑
			 */

			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Task %v processed\n", taskID)
			return nil
		}, i)

		// 所有的任务放入list中
		allTask = append(allTask, task)
	}

	// 建立一个池，
	// input:待处理的任务对列;池数量
	pool := workerpool.NewPool(allTask, 5)
	//pool.Run()


	go func() {
		for {
			taskID := rand.Intn(100) + 20

			if taskID % 7 == 0 {
				klog.Info("taskID: ", taskID, "pool stop!")
				pool.Stop()
			}

			time.Sleep(time.Duration(rand.Intn(5))*time.Second)
			task := workerpool.NewTask(func(data interface{}) error {
				taskID := data.(int)
				time.Sleep(100*time.Millisecond)
				fmt.Printf("Task %v processed \n", taskID)
				return nil
			}, taskID)

			pool.AddTask(task)
		}



	}()

	pool.RunBackground()






}
