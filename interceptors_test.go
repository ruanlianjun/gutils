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
		return func(data *int64) error {
			*data -= 1
			fmt.Println("start-1", *data)
			err := next(data)
			if err != nil {
				panic(err)
			}
			return nil
		}
	}, func(next NextInterceptorFunction[*int64]) NextInterceptorFunction[*int64] {
		return func(data *int64) error {
			*data -= 1
			fmt.Println("start-2", *data)
			err := next(data)
			if err != nil {
				panic(err)
			}

			return nil
		}
	})
	if err != nil {
		panic(err)
	}
}
