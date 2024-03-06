package main

import "github.com/myconcurrencytools/workpoolframework/example/http_example/cmd"

// 启动：
// ➜  http_example git:(main) ✗ go run main.go httpServer -p=8081
// 2024/03/06 12:18:14 start scheduler...
// I0306 12:18:14.325710    8727 worker.go:102] Starting worker background: 6
// I0306 12:18:14.325697    8727 worker.go:102] Starting worker background: 4
// I0306 12:18:14.325724    8727 worker.go:102] Starting worker background: 5
// I0306 12:18:14.325752    8727 pool.go:142] no task in global queue...
// I0306 12:18:14.325765    8727 worker.go:102] Starting worker background: 1
// I0306 12:18:14.325752    8727 worker.go:102] Starting worker background: 3
// I0306 12:18:14.325830    8727 worker.go:102] Starting worker background: 2

func main() {
	cmd.Execute()
}
