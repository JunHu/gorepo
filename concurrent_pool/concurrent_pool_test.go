package concurrent_pool

import (
	"fmt"
	"testing"
)

func Test_ConcurrentPool(t *testing.T) {
	tasks := make([]*Task, 100)
	for i := 0; i < 100; i += 1 {
		n := i
		Func := func() (interface{}, error) {
			return n * 2, nil
		}
		task := NewTask(i, "test_func", Func)
		tasks[i] = task
	}
	pool := NewConcurrentPool(10, tasks)
	pool.Start()

	results := pool.Results()
	fmt.Println(results)
}
