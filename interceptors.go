package gutils

type NextInterceptorFunction[T any] func(data T) error

type InterceptorFunction[T any] func(function NextInterceptorFunction[T]) NextInterceptorFunction[T]

func RunInterceptor[T any](data T, next NextInterceptorFunction[T], interceptor ...InterceptorFunction[T]) error {
	l := len(interceptor)
	for i := l - 1; i >= 0; i-- {
		next = interceptor[i](next)
	}
	return next(data)
}
