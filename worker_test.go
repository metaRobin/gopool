package gopool

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type wtask struct {
	ctx context.Context
	i   int
}

func TestWorker(t *testing.T) {
	executor := func(t interface{}) error {
		t, ok := t.(*wtask)
		if !ok {
			fmt.Println("not expect task struct")
		}
		fmt.Println("do task:", t)
		return nil
	}

	bWorker := &Worker{
		config:    &PoolConfig{chanSize: 5, batchExe: nil},
		taskQueue: make(chan interface{}, 5),
		executor:  executor,
		finshed:   make(chan struct{}),
		id:        0,
	}
	go bWorker.start()

	for i := 0; i < 20; i++ {
		bWorker.addTask(&wtask{
			ctx: context.Background(),
			i:   i,
		})
		if i%7 == 0 {
			fmt.Println("time sleep...")
			time.Sleep(15 * time.Millisecond)
		}
	}
	bWorker.quit()
	t.Log("work finished")
}
