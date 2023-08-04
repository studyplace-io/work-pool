package workerpool

import "time"

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