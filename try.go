package celerceptor

import (
	"context"
	"time"
)

const MaxAttempts = 3

/*
retry: wie oft/lang, wenn immer selber err=>perm err sonst temp err
timeouts?
context+cancel
logging?
*/

type RetryableFunc func() error

type RetryableFuncF[T any] func() (T, error)

func Try(ctx context.Context, f RetryableFunc) error {
	_, err := TryF(
		ctx,
		func() (any, error) {
			return nil, f()
		},
	)

	return err
}

func TryF[T any](ctx context.Context, f RetryableFuncF[T]) (T, error) {
	var (
		currentDelay = time.Second * 3
		result       T
		err          error //?
	)

	for trys := 1; trys <= MaxAttempts; trys++ {
		if ctx.Err() != nil {
			return result, ctx.Err()
		}

		result, err = f()
		if err == nil {
			return result, err
		}

		select {
		case <-ctx.Done():
			return result, ctx.Err()

		case <-time.After(currentDelay):
			currentDelay = currentDelay * 2
			// retry in loop
		}
	}

	return result, err
}
