package gopool

// var batchTaskPool sync.Pool

// func getTaskSlice() []Task {
// 	// return batchTaskPool.Get().([]Task)
// 	return make([]interface{}, 0, 10)
// }

type Close interface {
	GraceQuit()
}

type Pool interface {
	AddTask(task Task)
	Close
}

type ConsistancyPool interface {
	AddHashTask(hash uint64, task Task)
	Close
}
