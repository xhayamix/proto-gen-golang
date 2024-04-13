package combination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testItem struct {
	value  int32
	weight int32
}

func (ti *testItem) GetValue() int32 {
	return ti.value
}

func (ti *testItem) GetWeight() int32 {
	return ti.weight
}

func TestSearch(t *testing.T) {
	res := Search(Items{
		&testItem{value: 240, weight: 5},
		&testItem{value: 60, weight: 10},
		&testItem{value: 100, weight: 20},
		&testItem{value: 120, weight: 30},
		&testItem{value: 120, weight: 40},
		&testItem{value: 10, weight: 50},
	}, 50, 180, 240)

	assert.Equal(t, []Items{
		{
			&testItem{value: 120, weight: 40},
			&testItem{value: 60, weight: 10},
		},
		{
			&testItem{value: 120, weight: 30},
			&testItem{value: 100, weight: 20},
		},
		{
			&testItem{value: 120, weight: 30},
			&testItem{value: 60, weight: 10},
		},
		{
			&testItem{value: 240, weight: 5},
		},
	}, res)
}
