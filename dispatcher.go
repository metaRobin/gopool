package gopool

import (
	"fmt"
	"sync"
	"sync/atomic"
)

const (
	dispatcherStatusNone = 0
	dispatcherStatusRunning = 1
	dispatcherStatusStopped = 2
)

// Dispatcher 一致性派发器
type Dispatcher struct {
	config         *PoolConfig
	wg             *sync.WaitGroup
	quit           chan struct{}
	started	       chan struct{}
	status         uint32 // 派发器状态
	batchExe BatchExecutor
	exe      Executor
	workers []Worker
	batchers []Batcher
	isBatch bool // 是否位批量处理
}

func NewDispatcher(options ...Option) *Dispatcher {
	config := &PoolConfig{
		batchSize:  10,
		workers:    2,
	}

	for i := range options {
		options[i](config)
	}

	if config.chanSize <= config.workers {
		config.chanSize = config.workers
	}

	c := &Dispatcher{
		config:   config,
		wg:       new(sync.WaitGroup),
		quit:     make(chan struct{}),
		started:  make(chan struct{}),
		status:   dispatcherStatusNone,
		exe: config.exe,
		batchExe: config.batchExe,
	}

	if c.exe == nil && c.batchExe == nil {
		panic("not set exe or batch exe")
	}

	if c.batchExe != nil {
		c.isBatch = true
	}

	go c.start()

	atomic.StoreUint32(&c.status, dispatcherStatusRunning)
	<-c.started // wait dispather start ok

	return c
}

func (c *Dispatcher) start() {
	if !c.isBatch {
		for i := 0; i < c.config.workers; i ++ {
			c.workers = append(c.workers, Worker{
				config:    c.config,
				taskQueue: make(chan interface{}, c.config.chanSize),
				executor:  c.exe,
				finshed:   make(chan struct{}),
				id:        i,
			})
		}
		
		for i := range c.workers {
			c.wg.Add(1)
			go c.workers[i].start()
		}
	} else {
		for i := 0; i < c.config.workers; i ++ {
			c.batchers = append(c.batchers, Batcher{
				config:    c.config,
				taskQueue: make(chan interface{}, c.config.chanSize),
				executor:  c.batchExe,
				id:        i,
			})
		}
		
		for i := range c.batchers {
			c.wg.Add(1)
			go c.batchers[i].start()
		}
	}
	

	c.started <- struct{}{}

	for {
		select {
		case <-c.quit:
			if c.isBatch {
				for i := range c.batchers {
					c.batchers[i].quit()
					c.wg.Done()
				}
			} else {
				for i := range c.workers {
					c.workers[i].quit()
					c.wg.Done()
				}
			}
		}
	}
}

func (c *Dispatcher) AddHashTask(h uint64, task Task) error {
	if c.status != dispatcherStatusRunning {
		return fmt.Errorf("hash dispatcher is stopped")
	}
	
	if !c.isBatch {
		idx := h % uint64(len(c.workers))
		return c.workers[idx].addTask(task)
	}

	if c.isBatch {
		idx := h % uint64(len(c.batchers))
		return c.batchers[idx].addTask(task)
	}
	
	return fmt.Errorf("unknow executor to work")
}

func (c *Dispatcher) GraceQuit() {
	atomic.StoreUint32(&c.status, dispatcherStatusStopped)
	c.quit <- struct{}{}
	c.wg.Wait()
	return
}