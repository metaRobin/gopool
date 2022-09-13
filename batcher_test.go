package gopool

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type task struct {
	ctx context.Context
	i int
}

func TestBatcher(t *testing.T) {
	executor := func(tasks []interface{}) error {
		for i := range tasks {
			t, ok := tasks[i].(*task)
			if !ok {
				fmt.Println("not expect task struct")
				continue
			}
			fmt.Println("do task:", t.i)
		}
		fmt.Println("tasks:", tasks)
		return nil
	}
	
	bWorker := &Batcher{
		config:    &PoolConfig{
			batchSize:  5,
			chanSize:   5,
			linger: 10*time.Millisecond,
			batchErrCb: nil,
		},
		taskQueue: make(chan interface{}, 5),
		executor: executor,
	}
	go bWorker.start()

	for i := 0; i < 20; i ++ {
		bWorker.addTask(&task{
			ctx: context.Background(),
			i :i,
		})
		if i%7==0 {
			fmt.Println("time sleep...")
			time.Sleep(15*time.Millisecond)
		}
	}
	bWorker.quit()
	time.Sleep(11*time.Millisecond)
	t.Log("batcher finished")
}