package structures_test

import (
	"testing"
	"time"

	"github.com/stevo-go-utils/structures"
)

func TestConcurrency(t *testing.T) {
	connHandler := structures.NewConcurrencyHandler(2)
	connHandler.Start()
	ranTaskIdxes := []int{}
	for i := 0; i < 100; i++ {
		go func(i int) {
			task := connHandler.Enqueue(func() {
				ranTaskIdxes = append(ranTaskIdxes, i)
				time.Sleep(time.Millisecond * 100)
			})
			if i%2 == 0 {
				task.Cancel()
			}
		}(i)
	}
	connHandler.Done()
	connHandler.Wait()
}
