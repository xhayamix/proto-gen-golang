package parallel

import (
	"context"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"

	"github.com/xhayamix/proto-gen-golang/protoc-gen-all/perrors"
)

const DefaultSize int32 = 100

type Group struct {
	eg  *errgroup.Group
	sem *semaphore.Weighted
}

func NewGroupWithContext(ctx context.Context, size int32) (*Group, context.Context) {
	eg, egctx := errgroup.WithContext(ctx)
	return &Group{
		eg:  eg,
		sem: semaphore.NewWeighted(int64(size)),
	}, egctx
}

func (g *Group) Go(ctx context.Context, f func(ctx context.Context) error) {
	if g.eg == nil || g.sem == nil {
		return
	}

	g.eg.Go(func() error {
		if err := g.sem.Acquire(ctx, 1); err != nil {
			return perrors.Stack(err)
		}
		defer g.sem.Release(1)

		if err := f(ctx); err != nil {
			return perrors.Stack(err)
		}

		return nil
	})
}

func (g *Group) Wait() error {
	if g.eg == nil {
		return nil
	}

	if err := g.eg.Wait(); err != nil {
		return perrors.Stack(err)
	}

	return nil
}
