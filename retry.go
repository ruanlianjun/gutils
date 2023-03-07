package gutils

const (
	defaultRetry = 3
)

type (
	retryOptions struct {
		times int
	}
	RetryOption func(*retryOptions)
)

func DoWithRetry(fn func() error, opts ...RetryOption) error {
	options := newRetryOptions()
	for _, opt := range opts {
		opt(options)
	}
	var batchError *BatchError
	for i := 0; i < options.times; i++ {
		if err := fn(); err == nil {
			return nil
		} else {
			batchError.Add(err)
		}
	}
	return batchError.Err()
}

func newRetryOptions() *retryOptions {
	return &retryOptions{times: defaultRetry}
}
