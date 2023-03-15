package gutils

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func TestRunInterceptor(t *testing.T) {
	var tmp int64 = 10
	data := atomic.LoadInt64(&tmp)
	err := RunInterceptor[*int64](&data, func(data *int64) error {
		fmt.Println("start", *data)
		return nil
	}, func(next NextInterceptorFunction[*int64]) NextInterceptorFunction[*int64] {
		data -= 1
		return func(data *int64) error {
			fmt.Println("start-1", *data)
			next(data)
			fmt.Println("after")
			return nil
		}
	}, func(next NextInterceptorFunction[*int64]) NextInterceptorFunction[*int64] {
		data -= 1
		return func(data *int64) error {
			fmt.Println("start-2", *data)
			next(data)
			fmt.Println("after-2")
			return nil
		}
	})
	if err != nil {
		panic(err)
	}
}
