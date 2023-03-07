package gutils

import (
	"sync"
)

const (
	defaultWorkers = 5
)

type (
	rxOptions struct {
		unlimitedWorkers bool
		workers          int
	}
	Option      func(opts *rxOptions)
	GenerateFun func(source chan<- any)
	ForAllFun   func(pipe <-chan any)
	ForEachFun  func(item any)
	WalkFunc    func(item any, pipe chan<- any)

	Stream struct {
		source <-chan any
	}
)

func From(generate GenerateFun) Stream {
	source := make(chan any)
	GoSafe(func() {
		defer close(source)
		generate(source)
	})
	return Range(source)
}

func Just(items ...any) Stream {
	source := make(chan any, len(items))

	for item := range items {
		source <- item
	}
	close(source)
	return Range(source)

}

func Range(source <-chan any) Stream {
	return Stream{
		source: source,
	}
}

func (s Stream) Buffer(n int) Stream {
	if n < 0 {
		n = 0
	}
	source := make(chan any, n)
	go func() {
		for item := range s.source {
			source <- item
		}
		close(source)
	}()
	return Range(source)
}

func (s Stream) ForAll(fn ForAllFun) {
	fn(s.source)
	go drain(s.source)
}

func (s Stream) ForEach(fn ForEachFun) {
	for item := range s.source {
		fn(item)
	}
}

// Split 将数据分成多少份返回
func (s Stream) Split(n int) Stream {
	if n < 1 {
		panic("分割数据不应该小于1")
	}
	source := make(chan any)
	go func() {
		chunk := make([]any, 0, n)
		for item := range s.source {
			chunk = append(chunk, item)
			if len(chunk) == n {
				source <- chunk
				chunk = make([]any, 0, n)
			}
		}
		if chunk != nil {
			source <- chunk
		}
		close(source)
	}()
	return Range(source)
}

func (s Stream) Done() {
	drain(s.source)
}

func (s Stream) Walk(fn WalkFunc, opts ...Option) Stream {
	option := buildOptions(opts...)
	if option.unlimitedWorkers {
		return s.walkUnlimited(fn, option)
	}

	return s.walkWithLimit(fn, option)
}

func (s Stream) walkUnlimited(fn WalkFunc, option *rxOptions) Stream {
	pipe := make(chan any, option.workers)
	go func() {
		wg := sync.WaitGroup{}
		for item := range s.source {
			val := item
			wg.Add(1)
			GoSafe(func() {
				defer wg.Done()
				fn(val, pipe)
			})
		}
		wg.Wait()
		close(pipe)
	}()
	return Range(pipe)
}

// 限制coroutine的数量
func (s Stream) walkWithLimit(fn WalkFunc, option *rxOptions) Stream {
	pipe := make(chan any, option.workers)
	go func() {
		var wg sync.WaitGroup
		// 限制go coroutine的数量

		pool := make(chan struct{}, option.workers)
		for item := range s.source {
			wg.Add(1)
			val := item
			pool <- struct{}{}
			GoSafe(func() {
				defer func() {
					wg.Done()
					<-pool
				}()
				fn(val, pipe)
			})
		}
		wg.Wait()
		close(pipe)
	}()

	return Range(pipe)
}

func buildOptions(opts ...Option) *rxOptions {
	options := newOptions()
	for _, opt := range opts {
		opt(options)
	}

	return options
}

func WithUnlimitedWorkers() Option {
	return func(opts *rxOptions) {
		opts.unlimitedWorkers = true
	}
}

func WithWorkers(worker int) Option {
	return func(opts *rxOptions) {
		opts.workers = worker
	}
}

func newOptions() *rxOptions {
	return &rxOptions{
		workers: defaultWorkers,
	}
}

func drain(channel <-chan any) {
	for _ = range channel {
	}
}
