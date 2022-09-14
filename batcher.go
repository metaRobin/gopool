package gopool

import (
	"fmt"
	"time"
)

type Batcher struct {
	config    *PoolConfig
	taskQueue chan Task
	executor  BatchExecutor
	id        int
}

func (bat *Batcher) start() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("batcher panic", "r", r)
			bat.start()
		}
	}()
	batchSize := bat.config.batchSize
	batch := make([]Task, 0, batchSize)
	lingerTime := bat.config.linger
	ticker := time.NewTimer(lingerTime)
	if !ticker.Stop() {
		<-ticker.C
	}
	defer ticker.Stop()

	for {
		select {
		case task, ok := <-bat.taskQueue:
			if !ok {
				bat.process(batch)
				batch = make([]Task, 0, batchSize)
				return
			}

			batch = append(batch, task)
			size := len(batch)
			if size < batchSize {
				if size == 1 { // first element, tick reset
					ticker.Reset(lingerTime)
				}
				break
			}
			bat.process(batch)
			if !ticker.Stop() {
				<-ticker.C
			}
			batch = make([]Task, 0, batchSize)
		case <-ticker.C:
			if len(batch) == 0 {
				break
			}
			bat.process(batch)
			batch = make([]Task, 0, batchSize)
		}
	}
}

func (bat *Batcher) process(tasks []Task) {
	if err := bat.executor(tasks); err != nil {
		fmt.Println("process task err", "err", err)
		if bat.config.batchErrCb != nil {
			go bat.config.batchErrCb(err, tasks, bat.executor)
			return
		}
	}
}

func (bat *Batcher) addTask(t Task) error {
	bat.taskQueue <- t
	return nil
}

func (bat *Batcher) quit() {
	close(bat.taskQueue)
	batch := make([]Task, 0, bat.config.batchSize)
	for {
		select {
		case task, ok := <-bat.taskQueue:
			if !ok {
				fmt.Println("batcher close quit...", "id", bat.id, "batch", len(batch))
				bat.process(batch)
				return
			}
			batch = append(batch, task)
		}
	}
}
