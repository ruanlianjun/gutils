package gutils

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestBuffer(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		const N = 5
		var count int32
		var wait sync.WaitGroup
		wait.Add(1)
		From(func(source chan<- any) {
			ticker := time.NewTicker(10 * time.Millisecond)
			defer ticker.Stop()

			for i := 0; i < 2*N; i++ {
				select {
				case source <- i:
					atomic.AddInt32(&count, 1)
				case <-ticker.C:
					wait.Done()
					return
				}
			}
		}).Buffer(N).ForAll(func(pipe <-chan any) {
			wait.Wait()
			// why N+1, because take one more to wait for sending into the channel
			fmt.Println(atomic.LoadInt32(&count))
		})
	})

}

func TestWalk(t *testing.T) {
	runCheckedTest(t, func(t *testing.T) {
		From(func(source chan<- any) {
			for i := 0; i < 100; i++ {
				source <- i
			}
		}).
			Walk(func(item any, pipe chan<- any) {

				pipe <- item
			}, WithWorkers(3)).
			ForAll(func(pipe <-chan any) {
				for item := range pipe {
					fmt.Printf("each:%v\n", item)
					time.Sleep(time.Second * 1)
				}
			})

	})
}

func runCheckedTest(t *testing.T, fn func(t *testing.T)) {
	defer goleak.VerifyNone(t)
	fn(t)
}
