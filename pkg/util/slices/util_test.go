package slices

import (
	"strconv"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func Test_Chunk(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		assert.Equal(t, [][]int{
			{1, 2},
			{3, 4},
			{5},
		}, Chunk([]int{1, 2, 3, 4, 5}, 2))
	})

	t.Run("string", func(t *testing.T) {
		assert.Equal(t, [][]string{
			{"1", "2"},
			{"3", "4"},
			{"5"},
		}, Chunk([]string{"1", "2", "3", "4", "5"}, 2))
	})

	t.Run("defined type", func(t *testing.T) {
		type mySlice []string
		assert.Equal(t, []mySlice{
			{"1", "2"},
			{"3", "4"},
			{"5"},
		}, Chunk(mySlice{"1", "2", "3", "4", "5"}, 2))
	})

	t.Run("blank slice", func(t *testing.T) {
		type mySlice []string
		assert.Equal(t, []mySlice{}, Chunk(mySlice{}, 2))
	})

	t.Run("size is zero or less", func(t *testing.T) {
		assert.Equal(t, [][]int{{1, 2, 3, 4, 5}}, Chunk([]int{1, 2, 3, 4, 5}, 0))
	})
}

func Test_Collect(t *testing.T) {
	assert.Equal(
		t,
		[]string{"1", "2", "3"},
		Collect([]int{1, 2, 3}, strconv.Itoa),
	)

	assert.Equal(
		t,
		[]int{1, 2, 3},
		Collect([]string{"X", "XX", "XXX"}, utf8.RuneCountInString),
	)
}

func TestCopy(t *testing.T) {
	is1 := []int{1, 2, 3}
	is2 := Copy(is1)
	assert.Equal(t, is1, is2)

	is2 = append(is2, 4)
	assert.NotEqual(t, is1, is2)

	ss1 := []string{"1", "2", "3"}
	ss2 := Copy(ss1)
	assert.Equal(t, ss1, ss2)

	ss2 = append(ss2, "4")
	assert.NotEqual(t, ss1, ss2)
}

func Test_Diff(t *testing.T) {
	assert.Equal(t, []int{1, 2}, Diff([]int{1, 2, 3}, []int{3}))
	assert.Equal(t, []int{3, 1, 2}, Diff([]int{3, 1, 2, 3}, []int{3, 4, 5}))
	assert.Equal(t, []int{1, 2, 3, 2, 3}, Diff([]int{1, 2, 3, 2, 3, 4, 5, 3}, []int{3, 4, 5}))

	assert.Equal(t, []string{"1", "2"}, Diff([]string{"1", "2", "3"}, []string{"3"}))
}

func Test_Equal(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		assert.Equal(t, true, Equal(
			[]int{1, 2, 3},
			[]int{1, 2, 3},
		))
		assert.Equal(t, true, Equal(
			[]int{},
			[]int{},
		))
	})

	t.Run("異常", func(t *testing.T) {
		assert.Equal(t, false, Equal(
			[]int{1, 3, 2},
			[]int{1, 2, 3},
		))
		assert.Equal(t, false, Equal(
			[]int{1, 2, 3},
			[]int{1, 2},
		))
	})
}

func Test_Filter(t *testing.T) {
	assert.Equal(
		t,
		[]int{4, 5},
		Filter([]int{1, 2, 3, 4, 5}, func(n int) bool { return n > 3 }),
	)

	assert.Equal(
		t,
		[]string{"XX", "XXX"},
		Filter([]string{"X", "XX", "XXX"}, func(n string) bool { return utf8.RuneCountInString(n) > 1 }),
	)
}

func Test_FilterCollect(t *testing.T) {
	assert.Equal(
		t,
		[]string{"1", "2", "3"},
		FilterCollect([]int{1, 2, 3, 4, 5}, func(n int) (string, bool) { return strconv.Itoa(n), n < 4 }),
	)

	assert.Equal(
		t,
		[]int{3},
		FilterCollect([]string{"X", "XX", "XXX"}, func(n string) (int, bool) { return utf8.RuneCountInString(n), n == "XXX" }),
	)
}

func Test_First(t *testing.T) {
	assert.Equal(t, 4, First([]int{1, 2, 3, 4, 5}, func(n int) bool { return n > 3 }))
	assert.Equal(t, "XX", First([]string{"X", "XX", "XXX"}, func(n string) bool { return utf8.RuneCountInString(n) > 1 }))
}

func Test_Flatten(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, Flatten([][]int{{1, 2}, {3}}))
	assert.Equal(t, []string{"1", "2", "3"}, Flatten([][]string{{"1", "2"}, {"3"}}))
}

func Test_Has(t *testing.T) {
	assert.True(t, Has([]int{1, 2, 3, 4, 5}, func(n int) bool { return n == 3 }))
	assert.False(t, Has([]string{"X", "XX", "XXX"}, func(n string) bool { return utf8.RuneCountInString(n) > 3 }))
}

func Test_PartitionByIndex(t *testing.T) {
	genList := func(count int) []int {
		list := make([]int, 0, count)
		for i := 1; i <= count; i++ {
			list = append(list, i)
		}
		return list
	}

	t.Run("sizeが1以下", func(t *testing.T) {
		list := genList(3)
		assert.Equal(t, [][]int{
			{1, 2, 3},
		}, PartitionByIndex(list, 1))
	})

	t.Run("割り切れる", func(t *testing.T) {
		list := genList(9)
		assert.Equal(t, [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
		}, PartitionByIndex(list, 3))
	})

	t.Run("割り切れない", func(t *testing.T) {
		list := genList(11)
		assert.Equal(t, [][]int{
			{1, 2, 3},
			{4, 5, 6},
			{7, 8, 9},
			{10, 11},
		}, PartitionByIndex(list, 4))
	})

	t.Run("listがsize以下", func(t *testing.T) {
		list := genList(3)
		assert.Equal(t, [][]int{
			{1},
			{2},
			{3},
			{},
		}, PartitionByIndex(list, 4))
	})
}

func Test_Reverse(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		list := []int{1, 2, 3, 4, 5}
		res := Reverse(list)
		assert.Equal(t, res, []int{5, 4, 3, 2, 1})
		assert.Equal(t, list, []int{1, 2, 3, 4, 5})
	})

	t.Run("正常: 空", func(t *testing.T) {
		assert.Empty(t, Reverse([]int{}))
	})
}

func Test_Shuffle(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		list := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
		res := Shuffle(list)
		// rand.Shuffleの並べ替え方法についてはここでは考慮しない
		assert.ElementsMatch(t, res, list)
	})

	t.Run("正常: 空", func(t *testing.T) {
		assert.Empty(t, Shuffle([]int{}))
	})
}

func Test_Sort(t *testing.T) {
	t.Run("正常", func(t *testing.T) {
		list := []int{1, 2, 3, 4, 5}
		res := Sort(list, func(i, j int) bool { return i > j })
		assert.Equal(t, res, []int{5, 4, 3, 2, 1})
		assert.Equal(t, list, []int{1, 2, 3, 4, 5})
	})

	t.Run("正常: 空", func(t *testing.T) {
		assert.Empty(t, Sort([]int{}, func(i, j int) bool { return i > j }))
	})
}

func Test_Sum(t *testing.T) {
	assert.Equal(
		t,
		int(6),
		Sum([]int{1, 2, 3}, func(n int) int { return n }),
	)

	assert.Equal(
		t,
		int32(30),
		Sum([]int32{1, 2, 3}, func(n int32) int32 { return n * 5 }),
	)

	assert.Equal(
		t,
		float32(1.23),
		Sum([]float32{1.0, 0.2, 0.03}, func(n float32) float32 { return n }),
	)
}

func Test_Take(t *testing.T) {
	for name, tt := range map[string]struct {
		size     int
		expected []int
	}{
		"size: 0": {
			size:     0,
			expected: []int{},
		},
		"size: 1": {
			size:     1,
			expected: []int{1},
		},
		"size: 2": {
			size:     2,
			expected: []int{1, 2},
		},
		"size: 3": {
			size:     3,
			expected: []int{1, 2},
		},
	} {
		tt := tt
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, Take([]int{1, 2}, tt.size))
		})
	}
}

func Test_ToMap(t *testing.T) {
	assert.Equal(
		t,
		map[string]int{"1": 1, "2": 2, "3": 3},
		ToMap([]int{1, 2, 3}, strconv.Itoa),
	)

	assert.Equal(
		t,
		map[int]string{1: "X", 2: "XX", 3: "XXX"},
		ToMap([]string{"X", "XX", "XXX"}, utf8.RuneCountInString),
	)
}
