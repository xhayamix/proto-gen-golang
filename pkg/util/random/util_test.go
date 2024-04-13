package random

import (
	"testing"

	"github.com/scylladb/go-set/iset"
	"github.com/stretchr/testify/assert"
)

type drawable struct {
	index int
	ratio int64
}

func (d *drawable) GetRatio() int64 {
	return d.ratio
}

func TestDraw(t *testing.T) {
	rand := New(1, true)
	drawables := make([]Pickable, 0, 100)
	for i := 0; i < 100; i++ {
		drawables = append(drawables, &drawable{
			ratio: 1,
		})
	}
	for i := 0; i < 100; i++ {
		assert.NotNil(t, rand.Pick(drawables))
	}
}

func TestBulkDraw(t *testing.T) {
	rand := New(1, true)
	drawables := make([]Pickable, 0, 100)
	for i := 0; i < 100; i++ {
		drawables = append(drawables, &drawable{
			ratio: 1,
		})
	}
	for i := 0; i < 100; i++ {
		result := rand.BulkPick(5, drawables)
		assert.Equal(t, 5, len(result))
	}
}

func TestBulkDrawNoDuplication(t *testing.T) {
	rand := New(1, true)
	drawables := make([]Pickable, 0, 100)
	for i := 0; i < 100; i++ {
		drawables = append(drawables, &drawable{
			index: i,
			ratio: 1,
		})
	}
	count := 5
	for i := 0; i < 100; i++ {
		result := rand.BulkPickNoDuplication(count, drawables)
		assert.Equal(t, count, len(result))
		indexSet := iset.NewWithSize(len(result))
		for _, r := range result {
			indexSet.Add(r.(*drawable).index)
		}
		assert.Equal(t, count, indexSet.Size())
	}
}

func TestNew(t *testing.T) {
	rand := New(1, true)
	wants := []int{81, 887, 847, 59, 81, 318, 425, 540, 456, 300}
	value := 1000
	for _, want := range wants {
		assert.Equal(t, want, rand.NextIntn(value))
	}
}

func TestRandImpl_NextBool(t *testing.T) {
	rand := New(1, true)
	assert.False(t, rand.NextBool())
}

func TestRandImpl_NextIntn(t *testing.T) {
	rand := New(1, true)
	assert.Equal(t, 81, rand.NextIntn(1000))
}

func TestRandImpl_RangeIntn(t *testing.T) {
	rand := New(1, true)
	for i := 0; i < 100; i++ {
		res := rand.RangeIntn(100, 200)
		assert.True(t, res >= 100 && res <= 200)
	}
}

func TestRandImpl_NextInt31n(t *testing.T) {
	rand := New(1, true)
	assert.Equal(t, int32(81), rand.NextInt31n(1000))
}

func TestRandImpl_RangeInt31n(t *testing.T) {
	rand := New(1, true)
	for i := 0; i < 100; i++ {
		res := rand.RangeInt31n(100, 200)
		assert.True(t, res >= 100 && res <= 200)
	}
}

func TestRandImpl_NextInt63n(t *testing.T) {
	rand := New(1, true)
	assert.Equal(t, int64(410), rand.NextInt63n(1000))
}

func TestRandImpl_RangeInt63n(t *testing.T) {
	rand := New(1, true)
	for i := 0; i < 100; i++ {
		res := rand.RangeInt63n(100, 200)
		assert.True(t, res >= 100 && res <= 200)
	}
}

func TestRandImpl_Hit(t *testing.T) {
	rand := New(1, true)
	assert.True(t, rand.Hit(100, 1000))
	assert.False(t, rand.Hit(100, 1000))
	assert.False(t, rand.Hit(100, 1000))

	// 100%と0%の確認
	assert.True(t, rand.Hit(1000, 1000))
	assert.False(t, rand.Hit(0, 1000))
}

func TestRandImpl_HitPercent(t *testing.T) {
	rand := New(1, true)
	assert.False(t, rand.HitPercent(50))
	assert.False(t, rand.HitPercent(50))
	assert.True(t, rand.HitPercent(50))

	// 100%と0%の確認
	assert.True(t, rand.HitPercent(100))
	assert.False(t, rand.HitPercent(0))
}

func TestRandImpl_HitPermil(t *testing.T) {
	rand := New(1, true)
	assert.True(t, rand.HitPermil(500))
	assert.False(t, rand.HitPermil(500))
	assert.False(t, rand.HitPermil(500))

	// 100%と0%の確認
	assert.True(t, rand.HitPermil(1000))
	assert.False(t, rand.HitPermil(0))
}

func TestRandImpl_HitPermyriad(t *testing.T) {
	rand := New(1, true)
	assert.False(t, rand.HitPermyriad(5000))
	assert.False(t, rand.HitPermyriad(5000))
	assert.True(t, rand.HitPermyriad(5000))

	// 100%と0%の確認
	assert.True(t, rand.HitPermyriad(10000))
	assert.False(t, rand.HitPermyriad(0))
}

func TestRandImpl_Shuffle(t *testing.T) {
	rand := New(1, true)
	slice := []int{1, 2, 3, 4, 5, 6}
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
	assert.Equal(t, []int{6, 1, 2, 3, 5, 4}, slice)
}
