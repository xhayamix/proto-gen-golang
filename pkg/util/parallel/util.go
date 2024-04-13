package parallel

import (
	"context"
	"fmt"
	"runtime"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type Function func(ctx context.Context) error
type Functions []Function

func Run(ctx context.Context, functions Functions) error {
	eg, ctx := errgroup.WithContext(ctx)
	for _, f := range functions {
		eg.Go(wrapRecover(ctx, f))
	}
	return eg.Wait()
}

func recoverError(r interface{}) error {
	var stacktrace string
	for depth := 0; ; depth++ {
		_, file, line, ok := runtime.Caller(depth)
		if !ok {
			break
		}
		stacktrace += fmt.Sprintf("        %v:%d\n", file, line)
	}
	return fmt.Errorf("panic recovered: %v\n%s", r, stacktrace)
}

func wrapRecover(ctx context.Context, f Function) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverError(r)
			}
		}()
		return f(ctx)
	}
}

func RunConcurrency(ctx context.Context, functions Functions, concurrency int64) error {
	eg, ctx := errgroup.WithContext(ctx)
	sem := semaphore.NewWeighted(concurrency)
	for _, f := range functions {
		eg.Go(wrapRecoverConcurrency(ctx, sem, f))
	}
	return eg.Wait()
}

func wrapRecoverConcurrency(ctx context.Context, sem *semaphore.Weighted, f Function) func() error {
	return func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = recoverError(r)
			}
		}()
		if err := sem.Acquire(ctx, 1); err != nil {
			return err
		}
		defer sem.Release(1)
		return f(ctx)
	}
}
