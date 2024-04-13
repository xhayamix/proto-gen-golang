package parallel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Run(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		var res1, res2 int32
		functions := Functions{
			func(ctx context.Context) error {
				res1 = 1
				return nil
			},
			func(ctx context.Context) error {
				res2 = 2
				return nil
			},
		}
		err := Run(context.Background(), functions)

		assert.Equal(t, int32(1), res1)
		assert.Equal(t, int32(2), res2)
		assert.NoError(t, err)
	})

	t.Run("異常: panic recover", func(t *testing.T) {
		functions := Functions{
			func(ctx context.Context) error {
				panic("Panic!")
			},
		}

		assert.Error(t, Run(context.Background(), functions))
	})
}

func Test_RunConcurrency(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		var res1, res2 int32
		functions := Functions{
			func(ctx context.Context) error {
				res1 = 1
				return nil
			},
			func(ctx context.Context) error {
				res2 = 2
				return nil
			},
		}
		err := RunConcurrency(context.Background(), functions, 1)

		assert.Equal(t, int32(1), res1)
		assert.Equal(t, int32(2), res2)
		assert.NoError(t, err)
	})

	t.Run("正常: concurrency > len(functions)", func(t *testing.T) {
		var res1, res2 int32
		functions := Functions{
			func(ctx context.Context) error {
				res1 = 1
				return nil
			},
			func(ctx context.Context) error {
				res2 = 2
				return nil
			},
		}
		err := RunConcurrency(context.Background(), functions, 10)

		assert.Equal(t, int32(1), res1)
		assert.Equal(t, int32(2), res2)
		assert.NoError(t, err)
	})

	t.Run("異常: panic recover", func(t *testing.T) {
		functions := Functions{
			func(ctx context.Context) error {
				panic("Panic!")
			},
		}

		assert.Error(t, RunConcurrency(context.Background(), functions, 1))
	})
}
