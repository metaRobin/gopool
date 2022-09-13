package gopool

import "time"

type Task = interface{}  // 任务

// the func to batch process items
type BatchExecutor = func([]Task) error

// Executor 
type Executor = func(Task) error

// the func to handle error
type BatchErrorCallback = func(err error, items []Task, processor BatchExecutor)

// ErrorCallback error callback
type ErrorCallback = func(err error, item Task, processor Executor)

type Option = func(c *PoolConfig)

type PoolConfig struct {
	batchSize         int
	workers           int
	chanSize       int
	linger        time.Duration
	batchExe      BatchExecutor
	batchErrCb      BatchErrorCallback
	exe    Executor
	errCb  ErrorCallback
}

func BatchSize(sz int) Option {
	return func(c *PoolConfig) {
		c.batchSize = sz
	}
}

func WorkerCount(count int) Option {
	return func(c *PoolConfig) {
		c.workers = count
	}
}

func ChanSize(size int) Option {
	return func(c *PoolConfig) {
		c.chanSize = size
	}
}

// Linger linger time for batcher
func Linger(d time.Duration) Option {
	return func(c *PoolConfig) {
		c.linger = d
	}
}

func BatchExecute(batchExe BatchExecutor) Option {
	return func(c *PoolConfig) {
		c.batchExe = batchExe
	}
}

// BatchErrCallback batcher handle error callback
func BatchErrCallback(batchErrCb BatchErrorCallback) Option {
	return func(c *PoolConfig) {
		c.batchErrCb = batchErrCb
	}
}

func Execute(exe Executor) Option {
	return func(c *PoolConfig) {
		c.exe = exe
	}
}

// ErrCallback woker handle error callback
func ErrCallback(errCb ErrorCallback) Option {
	return func(c *PoolConfig) {
		c.errCb = errCb
	}
}