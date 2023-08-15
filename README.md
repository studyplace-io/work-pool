### work-pool
### 介绍
`work-pool`是基于golang实现的协程池，让调用者在使用并发时控制并发数量等功能，达到限制goroutine数量与复用的效果。

### 项目功能
- 自定义worker数量
- 自定义任务超时时间
- 自定义最大worker数，可根据task数自动扩缩容worker
- 自定义任务回调与错误回调方法(resultCallback、errorCallback)
- 支持阻塞式运行与非阻塞式运行

![项目架构](https://github.com/StudyPlace-io/work-pool-framework/blob/main/image/%E6%9E%B6%E6%9E%84%E6%96%B0%E5%9B%BE.png?raw=true)

### 使用
#### Pool配置
调用方可在初始化时决定Pool配置
- 超时时间
- 最大worker数
- 任务结束回调方法
- 错误处理回调方法
```go
// Option 选项模式
type Option func(pool *Pool)

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(p *Pool) {
		p.timeout = timeout
	}
}

// WithMaxWorkerNum 设置最大worker数量
func WithMaxWorkerNum(maxWorkerNum int) Option {
	return func(p *Pool) {
		p.maxWorkerNum = maxWorkerNum
	}
}

// WithResultCallback 设置结果回调方法
func WithResultCallback(callback func(interface{})) Option {
	return func(p *Pool) {
		p.resultCallback = callback
	}
}

// WithErrorCallback 设置错误回调方法
func WithErrorCallback(callback func(error)) Option {
	return func(p *Pool) {
		p.errorCallback = callback
	}
}
```
#### 基本使用
1. 实例化Pool
```go
 pool := workerpool.NewPool(5, workerpool.WithTimeout(1), workerpool.WithErrorCallback(func(err error) {
        fmt.Println("WithErrorCallback")
        if err != nil {
            panic(err)
        }
    }), workerpool.WithResultCallback(func(i interface{}) {
        fmt.Println("result: ", i)
    }))
```
2. 生成Pool可接受的任务
- 目前支持接口实现或使用内置的TaskInstance对象

调用方可实现此接口，即可视为Pool任务
```go
// Task 任务接口，由工作池抽象出的具体执行单元，
// 当pool启动时，会从chan中不断读取Task接口对象执行
type Task interface {
	Execute() (interface{}, error)
	GetTaskName() string
}
```
调用方处理好func逻辑后，也可直接使用内置的TaskInstance对象
```go
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
```
3. 放入池中
- 支持静态放入与动态放入

静态放入：Pool未启动时放入

`pool.AddGlobalQueue(task) // 所有的任务放入全局队列中`

动态放入：Pool启动时放入

`pool.AddTask(task)`

#### 示例1 阻塞式运行
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

#### 示例2 非阻塞式运行
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
    pool := workerpool.NewPool(5, workerpool.WithTimeout(1), workerpool.WithMaxWorkerNum(25),, workerpool.WithErrorCallback(func(err error) {
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

#### 更多示例：
可在/example目录下查看：
1. 封装简易调度器
2. 简易http服务实现执行任务