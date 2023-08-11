package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/studyplace-io/work-pool-framework/example/http_example/pkg/common"
	"github.com/studyplace-io/work-pool-framework/example/http_example/pkg/scheduler"
	"github.com/studyplace-io/work-pool-framework/example/http_example/pkg/server/model"
)

func HttpServer(c *common.ServerConfig) {

	if !c.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 启动调度器
	s := scheduler.NewScheduler(6)
	s.Start()

	r.GET("/test", func(c *gin.Context) {
		c.String(200, "测试用")
	})

	// 执行任务接口
	/*
	  url: http://127.0.0.1:8080/start
	  body入参：
	  {
	      "taskName": "tttt",
	      "taskType": "int",
	      "input": "aaaa"
	  }
	*/
	r.POST("/start", func(c *gin.Context) {
		var myTask model.MyTask
		if err := c.ShouldBindJSON(&myTask); err != nil {

			fmt.Errorf("do operation action error: %s", err)
			c.JSON(400, gin.H{"message": "start task error"})
			return
		}
		myTask.ChooseTaskType()
		s.AddTask(&myTask)

		model.TaskMap[myTask.TaskName] = &myTask
		c.JSON(200, gin.H{"message": "start task success"})
	})

	// 查询任务状态接口
	/*
	  http://127.0.0.1:8080/task?taskName=tttt
	*/
	r.GET("/task", func(c *gin.Context) {
		taskName := c.Query("taskName")
		my := model.TaskMap[taskName]
		c.JSON(200, gin.H{"message": my.Status})
	})

	err := r.Run(fmt.Sprintf(":%v", c.Port))
	fmt.Println(err)
}
