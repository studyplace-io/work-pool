# work-pool-framework
![项目架构](https://github.com/googs1025/Simple-work-pool-framework/blob/main/image/%E6%9E%B6%E6%9E%84.jpg?raw=true)

### 示例1 
**Run方法调用**
```go
/*
	使用方法：
	1. 创建工作池
	2. 定义需要的任务func
	3. 遍历任务数，放入全局队列
	4. 启动工作池
*/

func TestTaskPool1(t *testing.T) {

    pool := workerpool.NewPool(5, workerpool.WithTimeout(1), workerpool.WithErrorCallback(func(err error) {
        fmt.Println("WithErrorCallback")
        if err != nil {
            panic(err)
        }
    }), workerpool.WithResultCallback(func(i interface{}) {
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

```

### 示例2
```go
/*
	使用方法：
	1. 创建工作池
	2. 定义需要的任务func
	3. 遍历任务数，放入全局队列
	4. 启动工作池
*/

func TestTaskPool2(t *testing.T) {

    // 建立一个池，
    // input:池数量

    //pool := workerpool.NewPool(5)
    pool := workerpool.NewPool(5, workerpool.WithTimeout(1), workerpool.WithErrorCallback(func(err error) {
        if err != nil {
            panic(err)
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


```