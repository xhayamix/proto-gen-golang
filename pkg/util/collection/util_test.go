package collection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type user struct {
	name string
}

type users []*user

func TestToMap(t *testing.T) {
	assert.Equal(t, map[string]*user{"a": {name: "a"}, "b": {name: "b"}}, ToMap(users{{name: "a"}, {name: "b"}}, func(e *user) string {
		return e.name
	}))
}

func TestSliceToMap(t *testing.T) {
	slice := users{
		{name: "a"},
		{name: "a"},
		{name: "b"},
	}

	got := ToSliceMap(slice, func(e *user) string {
		return e.name
	})

	assert.Len(t, got, 2)
	assert.Equal(t, users{
		{name: "a"},
		{name: "a"},
	}, got["a"])
	assert.Equal(t, users{
		{name: "b"},
	}, got["b"])
}

func TestSelect(t *testing.T) {
	assert.Equal(t, []string{"a", "b"}, Select(users{{name: "a"}, {name: "b"}}, func(e *user) string {
		return e.name
	}))
}

func TestSort(t *testing.T) {
	assert.Equal(t, users{{name: "a"}, {name: "b"}}, Sort(users{{name: "a"}, {name: "b"}}, func(e1, e2 *user) bool {
		return e1.name < e2.name
	}))
}

func TestWhere(t *testing.T) {
	assert.Equal(t, users{{name: "a"}}, Where(users{{name: "a"}, {name: "b"}}, func(e *user) bool {
		return e.name == "a"
	}))
}

func TestAny(t *testing.T) {
	assert.True(t, Any([]bool{true, true}, func(b bool) bool { return b }))
	assert.True(t, Any([]bool{true, false}, func(b bool) bool { return b }))
	assert.False(t, Any([]bool{false, false}, func(b bool) bool { return b }))
}

func TestAll(t *testing.T) {
	assert.True(t, All([]bool{true, true}, func(b bool) bool { return b }))
	assert.False(t, All([]bool{true, false}, func(b bool) bool { return b }))
	assert.False(t, All([]bool{false, false}, func(b bool) bool { return b }))
}

func TestCopy(t *testing.T) {
	assert.Equal(t, []string{"a", "b"}, Copy([]string{"a", "b"}))
	assert.Equal(t, users{{name: "a"}, {name: "b"}}, Copy(users{{name: "a"}, {name: "b"}}))
}

func TestFirst(t *testing.T) {
	assert.Equal(t, "a", First([]string{"a", "b"}))
	assert.Equal(t, "", First([]string{}))

	assert.Equal(t, &user{name: "a"}, First(users{{name: "a"}, {name: "b"}}))
	assert.Nil(t, First(users{}))
}

func TestLast(t *testing.T) {
	assert.Equal(t, "b", Last([]string{"a", "b"}))
	assert.Equal(t, "", Last([]string{}))

	assert.Equal(t, &user{name: "b"}, Last(users{{name: "a"}, {name: "b"}}))
	assert.Nil(t, Last(users{}))
}

func TestSplit(t *testing.T) {
	assert.Equal(t, [][]string{{"a", "b"}, {"c"}}, Split([]string{"a", "b", "c"}, 2))
	assert.Equal(t, []users{{{name: "a"}, {name: "b"}}, {{name: "c"}}}, Split(users{{name: "a"}, {name: "b"}, {name: "c"}}, 2))
}

func TestConcat(t *testing.T) {
	assert.Equal(t, []string{"a", "b", "c"}, Concat([]string{"a", "b"}, []string{"c"}))
	assert.Equal(t, users{{name: "a"}, {name: "b"}, {name: "c"}}, Concat(users{{name: "a"}, {name: "b"}}, users{{name: "c"}}))
}
