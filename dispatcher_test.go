package gopool

import (
	"context"
	"testing"
	"time"
)

type dispTask struct {
	ctx context.Context
	i   int
	h   uint64
}

func (d *dispTask) GetHash() uint64 {
	return d.h
}

func TestDispatcherBatcher(t *testing.T) {
	var batchExecutor = func(tasks []interface{}) error {
		t.Log("handle tasks,sz:", len(tasks))
		for i := range tasks {
			tk, ok := tasks[i].(*dispTask)
			if !ok {
				t.Log("not expect task struct", "task", tasks[i])
				continue
			}
			t.Log("do task:", tk.i, "h", tk.h)
		}
		return nil
	}
	dispatcher := NewDispatcher(BatchExecute(batchExecutor), ChanSize(10), BatchSize(10), WorkerCount(5))

	for i := 0; i < 20; i++ {
		dispatcher.AddHashTask(uint64(i), &dispTask{
			ctx: context.Background(),
			i:   i,
			h:   uint64(i % 5),
		})
	}

	dispatcher.GraceQuit()
	t.Log("DispatcherBatcher done...")
}

func TestDispatcherLinger(t *testing.T) {
	var batchExecutor = func(tasks []interface{}) error {
		t.Log("handle tasks,sz:", len(tasks), "ts", time.Now().UnixMilli())
		for i := range tasks {
			tk, ok := tasks[i].(*dispTask)
			if !ok {
				t.Log("not expect task struct", "task", tasks[i])
				continue
			}
			t.Log("do task:", tk.i, "h", tk.h)
		}
		return nil
	}

	dispatcher := NewDispatcher(BatchExecute(batchExecutor), ChanSize(10),
		Linger(5*time.Millisecond), BatchSize(10), WorkerCount(5))

	for i := 0; i < 13; i++ {
		var dt = &dispTask{
			ctx: context.Background(),
			i:   i,
			h:   uint64(i % 5),
		}

		dispatcher.AddHashTask(dt.GetHash(), dt)
		if i > 10 {
			time.Sleep(10 * time.Millisecond)
			dispatcher.AddHashTask(dt.GetHash(), dt)
		}
	}

	dispatcher.GraceQuit()
	t.Log("DispatcherLinger done...")
}

func TestDispatcherWorker(t *testing.T) {
	var exec = func(task interface{}) error {
		tk := task.(*dispTask)
		t.Log("do task:", "i", tk.i, "h", tk.h)
		return nil
	}
	dispatcher := NewDispatcher(Execute(exec), ChanSize(10), BatchSize(10), WorkerCount(5))

	for i := 0; i < 20; i++ {
		dispatcher.AddHashTask(uint64(i), &dispTask{
			ctx: context.Background(),
			i:   i,
			h:   uint64(i % 5),
		})
	}

	dispatcher.GraceQuit()
	t.Log("DispatcherWorker done...")
}
