package slices

import (
	"math/rand"
	"sort"
)

type number interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64
}

// アルファベット順

// Chunk 各要素がN個となるようスライスを分割する
func Chunk[T any, S ~[]T](s S, size int) []S {
	if size <= 0 {
		return []S{s}
	}

	length := len(s)
	result := make([]S, 0, length/size+1)
	for i := 0; i < length; i += size {
		end := i + size
		if length < end {
			end = length
		}
		result = append(result, s[i:end])
	}

	return result
}

// Collect 引数の関数の返り値を集めた配列を返す (他言語で言うmap)
func Collect[T any, S ~[]T, Elm any](s S, f func(e T) Elm) []Elm {
	ret := make([]Elm, 0, len(s))
	for _, e := range s {
		ret = append(ret, f(e))
	}
	return ret
}

// Copy sliceのコピー
func Copy[T any, S ~[]T](s S) S {
	ns := make(S, len(s))
	copy(ns, s)
	return ns
}

// Diff 第一引数のSliceから第二引数のSliceを後方から引いた結果を返す
// 第二引数のみに存在する要素は無視する
func Diff[T comparable, S ~[]T](s1, s2 S) S {
	m := make(map[T]int, len(s1))
	for _, e := range s1 {
		m[e]++
	}
	for _, e := range s2 {
		m[e]--
	}
	ret := make(S, 0, len(s1))
	for _, e := range s1 {
		if m[e] <= 0 {
			continue
		}
		ret = append(ret, e)
		m[e]--
	}
	return ret
}

// Equal スライスの比較を行う
func Equal[E comparable](s1, s2 []E) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// Filter 引数の関数を満たす要素の配列を返す
func Filter[T any, S ~[]T](s S, f func(e T) bool) S {
	ret := make(S, 0, len(s))
	for _, e := range s {
		if f(e) {
			ret = append(ret, e)
		}
	}
	return ret
}

// FilterCollect 引数の関数を満たした返り値を集めた配列を返す
func FilterCollect[T any, S ~[]T, Elm any](s S, f func(e T) (Elm, bool)) []Elm {
	ret := make([]Elm, 0, len(s))
	for _, e := range s {
		if val, ok := f(e); ok {
			ret = append(ret, val)
		}
	}
	return ret
}

// First 引数の関数を満たす初めの要素を返す
func First[T any, S ~[]T](s S, f func(e T) bool) (ret T) {
	for _, e := range s {
		if f(e) {
			return e
		}
	}

	return ret
}

// Flatten 2次元Sliceを1次元Sliceにする
func Flatten[T any, S ~[]T](s []S) []T {
	if len(s) == 0 {
		return nil
	}

	ret := make([]T, 0, len(s[0])*len(s))
	for _, e := range s {
		ret = append(ret, e...)
	}
	return ret
}

// Has 引数の関数を満たす要素が存在するか
func Has[T any, S ~[]T](s S, f func(e T) bool) bool {
	for _, e := range s {
		if f(e) {
			return true
		}
	}
	return false
}

// Max 引数の関数を満たした返り値の最大を返す
func Max[T any, S ~[]T, Num number](s S, f func(e T) Num) Num {
	if len(s) == 0 {
		return 0
	}
	return f(Sort(s, func(i, j T) bool { return f(i) > f(j) })[0])
}

// Min 引数の関数を満たした返り値の最小を返す
func Min[T any, S ~[]T, Num number](s S, f func(e T) Num) Num {
	if len(s) == 0 {
		return 0
	}
	return f(Sort(s, func(i, j T) bool { return f(i) < f(j) })[0])
}

// PartitionByIndex スライスをN個の要素に分割する
// 余りが発生する場合は前詰めして要素を降っていく
// ex: PartitionByIndex([]int{1, 2, 3, 4, 5}, 3) -> [][]int{{1, 2}, {3, 4}, {5}}
func PartitionByIndex[T any, S ~[]T](s S, size int) []S {
	if size <= 1 {
		return []S{s}
	}

	result := make([]S, 0, size)
	div := len(s) / size
	mod := len(s) % size

	currentIndex := 0
	remainingModCount := mod
	for count := 0; count < size; count++ {
		if remainingModCount > 0 {
			result = append(result, s[currentIndex:currentIndex+div+1])
			currentIndex += div + 1
			remainingModCount--
		} else {
			result = append(result, s[currentIndex:currentIndex+div]) // 範囲外~範囲外のスライスを取ると空の配列になる
			currentIndex += div                                       // len(s) / size が < 0 の場合は丸められて0を加算し続ける
		}
	}

	return result
}

// Reverse 配列コピーした後に逆順にして返す
func Reverse[T any, S ~[]T](s S) S {
	ret := make(S, len(s))
	copy(ret, s)

	length := len(ret)
	half := length / 2

	for i := 0; i < half; i++ {
		j := length - 1 - i
		ret[i], ret[j] = ret[j], ret[i]
	}

	return ret
}

// Shuffle 配列コピーした後にmath/rand.Shuffleを用いた結果を返す
func Shuffle[T any, S ~[]T](s S) S {
	ret := make(S, len(s))
	copy(ret, s)

	rand.Shuffle(len(ret), func(i, j int) {
		ret[i], ret[j] = ret[j], ret[i]
	})
	return ret
}

// Sort 配列コピーした後にソートする
func Sort[T any, S ~[]T](s S, less func(i, j T) bool) S {
	ret := make(S, len(s))
	copy(ret, s)

	sort.Slice(ret, func(i, j int) bool { return less(ret[i], ret[j]) })
	return ret
}

// Sum 引数の関数を満たした返り値の集計を返す
func Sum[T any, S ~[]T, Num number](s S, f func(e T) Num) Num {
	var ret Num
	for _, e := range s {
		ret += f(e)
	}
	return ret
}

// Take 先頭からN個の要素を取得
func Take[T any, S ~[]T](s S, size int) S {
	if len(s) < size {
		return s
	}
	return s[:size]
}

// ToMap 引数の関数の返り値をKeyにしたmapを返す
func ToMap[T any, S ~[]T, K comparable](s S, f func(e T) K) map[K]T {
	ret := make(map[K]T, len(s))
	for _, e := range s {
		ret[f(e)] = e
	}
	return ret
}
