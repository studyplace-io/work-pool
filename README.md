# work-pool-framework
![项目架构](https://github.com/googs1025/Simple-work-pool-framework/blob/main/image/%E6%9E%B6%E6%9E%84.jpg?raw=true)

### 示例1 
**Run方法调用**
```go
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
```
### 示例2
```go
/*
	使用方法：
	1. 准备全局的任务队列，用于存放任务
	2. 定义需要的任务func
	3. 遍历任务数，放入全局队列
	4. 创建且启动工作池
*/

func TestTaskPool2(t *testing.T) {

        // 存放任务的全局队列
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
        // 启动在后台等待执行
        go pool.RunBackground()
        
        for {
            taskID := rand.Intn(100) + 20
            
            // 模拟一个退出条件
            if taskID%7 == 0 {
                klog.Info("taskID: ", taskID, "pool stop!")
                pool.StopBackground()
                break
            }
            
            time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
            // 模拟后续加入pool
            task := workerpool.NewTask(func(data interface{}) error {
                taskID := data.(int)
                time.Sleep(100 * time.Millisecond)
                klog.Info("Task ", taskID, " processed")
                return nil
            }, taskID)
            
            pool.AddTask(task)
        }
        
        fmt.Println("finished...")
}


```