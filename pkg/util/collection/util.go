package collection

import (
	"sort"
)

func ToMap[T any, V comparable](s []T, f func(T) V) map[V]T {
	ret := make(map[V]T, len(s))
	for _, e := range s {
		ret[f(e)] = e
	}
	return ret
}

func ToSliceMap[S ~[]T, T any, V comparable](s S, f func(T) V) map[V]S {
	sliceLenMap := make(map[V]int, len(s))
	for _, e := range s {
		sliceLenMap[f(e)]++
	}

	ret := make(map[V]S, len(sliceLenMap))
	for _, e := range s {
		key := f(e)

		slice, ok := ret[key]
		if !ok {
			slice = make(S, 0, sliceLenMap[key])
		}

		ret[key] = append(slice, e)
	}
	return ret
}

func Select[T, V any](s []T, f func(T) V) []V {
	ret := make([]V, 0, len(s))
	for _, e := range s {
		ret = append(ret, f(e))
	}
	return ret
}

func Sort[S ~[]T, T any](s S, f func(e1, e2 T) bool) S {
	ret := Copy(s)
	sort.Slice(ret, func(i, j int) bool {
		return f(ret[i], ret[j])
	})
	return ret
}

func Where[S ~[]T, T any](s S, f func(T) bool) S {
	ret := make(S, 0, len(s))
	for _, e := range s {
		if f(e) {
			ret = append(ret, e)
		}
	}
	return ret
}

func Any[S ~[]T, T any](s S, f func(T) bool) bool {
	for _, e := range s {
		if f(e) {
			return true
		}
	}
	return false
}

func All[S ~[]T, T any](s S, f func(T) bool) bool {
	for _, e := range s {
		if !f(e) {
			return false
		}
	}
	return true
}

func Copy[S ~[]T, T any](s S) S {
	ret := make(S, 0, len(s))
	ret = append(ret, s...)
	return ret
}

func First[T any](s []T) T {
	if len(s) == 0 {
		var v T
		return v
	}
	return s[0]
}

func Last[T any](s []T) T {
	if len(s) == 0 {
		var v T
		return v
	}
	return s[len(s)-1]
}

func Split[S ~[]T, T any](s S, size int) []S {
	length := len(s)
	splits := make([]S, 0, length/size+1)
	for i := 0; i < length; i += size {
		end := i + size
		if length < end {
			end = length
		}
		splits = append(splits, s[i:end])
	}
	return splits
}

func Concat[S ~[]T, T any](s1, s2 S) S {
	ret := make(S, 0, len(s1)+len(s2))
	ret = append(ret, s1...)
	ret = append(ret, s2...)
	return ret
}
