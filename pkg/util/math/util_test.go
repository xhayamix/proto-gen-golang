package math

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	t.Run("正常: int", func(t *testing.T) {
		assert.Equal(t, 1, Abs(1))
		assert.Equal(t, 0, Abs(0))
		assert.Equal(t, 1, Abs(-1))
	})
	t.Run("正常: int32", func(t *testing.T) {
		assert.Equal(t, int32(1), Abs(int32(1)))
		assert.Equal(t, int32(0), Abs(int32(0)))
		assert.Equal(t, int32(1), Abs(int32(-1)))
	})
	t.Run("正常: int64", func(t *testing.T) {
		assert.Equal(t, int64(1), Abs(int64(1)))
		assert.Equal(t, int64(0), Abs(int64(0)))
		assert.Equal(t, int64(1), Abs(int64(-1)))
	})
	t.Run("正常: float32", func(t *testing.T) {
		assert.Equal(t, float32(1.1), Abs(float32(1.1)))
		assert.Equal(t, float32(0), Abs(float32(0)))
		assert.Equal(t, float32(1.1), Abs(float32(-1.1)))
	})
	t.Run("正常: float64", func(t *testing.T) {
		assert.Equal(t, 1.1, Abs(1.1))
		assert.Equal(t, float64(0), Abs(float64(0)))
		assert.Equal(t, 1.1, Abs(-1.1))
	})
}

func TestRoundTenPermyriad(t *testing.T) {
	t.Run("切り捨てされる", func(t *testing.T) {
		assert.Equal(t, int64(33333), RoundTenPermyriad(1, 3))
	})
	t.Run("切り上げされる", func(t *testing.T) {
		assert.Equal(t, int64(66667), RoundTenPermyriad(2, 3))
	})
}

func Test_Sum(t *testing.T) {
	// int32
	assert.Equal(t, int32(10), Sum([]int32{1, 2, 3, 4}...))
	assert.Equal(t, int32(10), Sum[int32](5, 5))

	// int64
	assert.Equal(t, int64(10), Sum([]int64{1, 2, 3, 4}...))
	assert.Equal(t, int64(10), Sum[int64](5, 5))
}

func TestPermutation(t *testing.T) {
	assert.Equal(t, int64(60), Permutation(5, 3).Int64())
}

func TestFactorial(t *testing.T) {
	assert.Equal(t, int64(120), Factorial(5).Int64())
}

func TestCombination(t *testing.T) {
	assert.Equal(t, int64(10), Combination(5, 3).Int64())
}

func TestHomogeneous(t *testing.T) {
	assert.Equal(t, int64(35), Homogeneous(5, 3).Int64())
}

func TestSafeAddInt64(t *testing.T) {
	t.Run("正常: オーバーフローしない加算", func(t *testing.T) {
		assert.Equal(t, int64(10), SafeAddInt64(5, 5, math.MaxInt64, math.MinInt64))
	})
	t.Run("正常: オーバーフローする加算（プラス、最小）", func(t *testing.T) {
		assert.Equal(t, int64(math.MaxInt64), SafeAddInt64(math.MaxInt64, 1, math.MaxInt64, math.MinInt64))
	})
	t.Run("正常: オーバーフローする加算（プラス、最大）", func(t *testing.T) {
		assert.Equal(t, int64(math.MaxInt64), SafeAddInt64(math.MaxInt64, math.MaxInt64, math.MaxInt64, math.MinInt64))
	})
	t.Run("正常: オーバーフローする加算（マイナス、最小）", func(t *testing.T) {
		assert.Equal(t, int64(math.MinInt64), SafeAddInt64(math.MinInt64, math.MinInt64, math.MaxInt64, math.MinInt64))
	})
	t.Run("正常: オーバーフローする加算（マイナス、最大）", func(t *testing.T) {
		assert.Equal(t, int64(math.MinInt64), SafeAddInt64(math.MinInt64, -1, math.MaxInt64, math.MinInt64))
	})
	t.Run("正常: 最大値を超える加算", func(t *testing.T) {
		assert.Equal(t, int64(10), SafeAddInt64(0, 100, 10, math.MinInt64))
	})
	t.Run("正常: 最小値を下回る加算", func(t *testing.T) {
		assert.Equal(t, int64(-10), SafeAddInt64(0, -100, math.MaxInt64, -10))
	})
}
