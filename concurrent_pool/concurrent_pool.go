package concurrent_pool

import (
	"fmt"
	"sort"
)

type Task struct {
	f      func() (interface{}, error)
	idx    int
	name   string
	Err    error
	Result interface{}
}

type TaskSlice []*Task

func (s TaskSlice) Len() int           { return len(s) }
func (s TaskSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s TaskSlice) Less(i, j int) bool { return s[i].idx < s[j].idx }

func NewTask(idx int, name string, f func() (interface{}, error)) *Task {
	return &Task{
		idx:  idx,
		name: name,
		f:    f,
	}
}

func (task *Task) do() {
	defer func() {
		if err := recover(); err != nil {
			task.Err = fmt.Errorf("process task %s err %v", task.name, err)
		}
	}()
	result, err := task.f()
	task.Result = result
	task.Err = err
}

type ConcurrentPool struct {
	taskSize   int
	poolSize   int
	taskChan   chan *Task
	resultChan chan *Task
}

func NewConcurrentPool(size int, tasks []*Task) *ConcurrentPool {
	taskChan := make(chan *Task, len(tasks))
	resultChan := make(chan *Task, len(tasks))

	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)
	pool := &ConcurrentPool{taskChan: taskChan, poolSize: size, resultChan: resultChan, taskSize: len(tasks)}

	return pool
}

func (pool *ConcurrentPool) Start() {
	for i := 0; i < pool.poolSize; i += 1 {
		go pool.start_one_worker()
	}
}

func (pool *ConcurrentPool) start_one_worker() {
	for task := range pool.taskChan {
		fmt.Println(task.idx)
		task.do()
		pool.resultChan <- task
	}
}

func (pool *ConcurrentPool) Results() []*Task {
	results := make([]*Task, pool.taskSize)
	for i := 0; i < pool.taskSize; i += 1 {
		results[i] = <-pool.resultChan
	}
	sort.Sort(TaskSlice(results))
	return results
}
