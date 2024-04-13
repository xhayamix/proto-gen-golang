package math

import (
	"math/big"
)

type numbers interface {
	int | int32 | int64 | float32 | float64
}

func Abs[T numbers](num T) T {
	if num < 0 {
		return -num
	}
	return num
}

func Sum[N numbers](l ...N) N {
	var ret N
	for _, n := range l {
		ret += n
	}
	return ret
}

// Permutation 順列の計算
func Permutation(n, k int) *big.Int {
	v := big.NewInt(1)
	if 0 < k && k <= n {
		for i := 0; i < k; i++ {
			v.Mul(v, big.NewInt(int64(n-i)))
		}
	} else if k > n {
		v = big.NewInt(0)
	}
	return v
}

// Factorial 階乗の計算
func Factorial(n int) *big.Int {
	return Permutation(n, n-1)
}

// Combination 組み合わせ数の計算
func Combination(n, k int) *big.Int {
	return new(big.Int).Div(Permutation(n, k), Factorial(k))
}

// Homogeneous 重複あり組み合わせ数の計算
func Homogeneous(n, k int) *big.Int {
	return Combination(n+k-1, k)
}

// SafeAddInt64 安全なint64の加算（丸め込み）
func SafeAddInt64(a, b, max, min int64) int64 {
	// オーバーフローの場合の決定論的な挙動を利用してオーバーフローを検知する（https://go.dev/ref/spec#Arithmetic_operators）
	sum := a + b
	if (sum < a) == (b > 0) { // オーバーフローする場合
		if b > 0 {
			return max
		} else {
			return min
		}
	} else {
		if sum > max {
			return max
		}
		if sum < min {
			return min
		}
	}
	return a + b
}
