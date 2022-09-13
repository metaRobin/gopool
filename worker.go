package gopool

type Worker struct {
	config    *PoolConfig
	taskQueue chan Task
	executor  Executor
	finshed   chan struct{}
	id        int
}

func (w *Worker) start() {
	defer func() {
		if r := recover(); r != nil {
			w.start()
		}
	}()

	for task := range w.taskQueue {
		w.process(task)
	}

	w.finshed <- struct{}{}
}

func (w *Worker) process(task Task) {
	if err := w.executor(task); err != nil {
		if w.config.errCb != nil {
			go w.config.errCb(err, task, w.executor)
			return
		}
	}
}

func (w *Worker) addTask(t Task) error {
	w.taskQueue <- t
	return nil
}

func (w *Worker) quit() {
	close(w.taskQueue)
	<-w.finshed
}
