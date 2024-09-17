//go:generate mockgen -source=$GOFILE -destination=mock_$GOPACKAGE/mock_$GOFILE
//go:generate goimports -w --local "ç" mock_$GOPACKAGE/mock_$GOFILE
package random

import (
	"math/rand"

	"github.com/google/uuid"
)

const (
	tinyInt = 2
)

type Pickable interface {
	// GetRatio 重みを取得する
	GetRatio() int64
}

type Rand interface {
	// NextBool boolの値が戻る
	NextBool() bool
	// NextInt31n 0 <= result < value の値が戻る
	NextInt31n(n int32) int32
	// RangeInt31n s <= result <= e の値が戻る
	RangeInt31n(start int32, end int32) int32
	// NextInt63n 0 <= result < value の値が戻る
	NextInt63n(n int64) int64
	// RangeInt63n s <= result <= e の値が戻る
	RangeInt63n(start int64, end int64) int64
	// NextIntn 0 <= result < value の値が戻る
	NextIntn(n int) int
	// RangeIntn s <= result <= e の値が戻る
	RangeIntn(start int, end int) int
	// Hit rate > 0以上、denominator未満の乱数 のときtrue
	Hit(rate int64, denominator int32) bool
	// HitPercent rate > 0以上、100未満の乱数 のときtrue
	HitPercent(rate int64) bool
	// HitPermil rate > 0以上、1000未満の乱数 のときtrue
	HitPermil(rate int64) bool
	// HitPermyriad rate > 0以上、10000未満の乱数 のときtrue
	HitPermyriad(rate int64) bool
	// Pick 抽選する
	Pick(pickables []Pickable) Pickable
	// BulkPick 複数回抽選する（重複あり）
	BulkPick(drawCount int, pickables []Pickable) []Pickable
	// BulkPickNoDuplication 複数回抽選する（重複なし）
	BulkPickNoDuplication(drawCount int, origin []Pickable) []Pickable
	// NewRandomUUID UUIDを生成する
	NewRandomUUID() (string, error)
	// Shuffle シャッフルする
	Shuffle(n int, swap func(i int, j int))
}

type randImpl struct {
	*rand.Rand
}

func New(seed int64, isThreadSafe bool) Rand {
	source := rand.NewSource(seed)
	if isThreadSafe {
		source = &lockedSource{src: source}
	}
	return &randImpl{
		//nolint:gosec // weak random
		Rand: rand.New(source),
	}
}

func (r *randImpl) NextBool() bool {
	return r.NextIntn(tinyInt) == 0
}

func (r *randImpl) NextInt31n(n int32) int32 {
	return r.Int31n(n)
}

func (r *randImpl) RangeInt31n(s, e int32) int32 {
	return s + r.Int31n(e-s+1)
}

func (r *randImpl) NextInt63n(n int64) int64 {
	return r.Int63n(n)
}

func (r *randImpl) RangeInt63n(s, e int64) int64 {
	return s + r.Int63n(e-s+1)
}

func (r *randImpl) NextIntn(n int) int {
	return r.Intn(n)
}

func (r *randImpl) RangeIntn(s, e int) int {
	return s + r.Intn(e-s+1)
}

func (r *randImpl) Hit(rate int64, denominator int32) bool {
	return rate > int64(r.NextInt31n(denominator))
}

func (r *randImpl) HitPercent(rate int64) bool {
	return r.Hit(rate, 100)
}

func (r *randImpl) HitPermil(rate int64) bool {
	return r.Hit(rate, 1000)
}

func (r *randImpl) HitPermyriad(rate int64) bool {
	return r.Hit(rate, 10000)
}

func (r *randImpl) Pick(pickables []Pickable) Pickable {
	// ratioにマイナスが設定されると抽選ロジックがバグるので念のため除外する
	target := make([]Pickable, 0, len(pickables))
	for _, pickable := range pickables {
		if pickable.GetRatio() > 0 {
			target = append(target, pickable)
		}
	}
	if len(target) == 0 {
		return nil
	}

	var totalWeight int64
	for _, pickable := range target {
		totalWeight += pickable.GetRatio()
	}
	if totalWeight <= 0 {
		return nil
	}
	random := r.NextInt63n(totalWeight)
	var temp int64
	for _, pickable := range target {
		temp += pickable.GetRatio()
		if temp > random {
			return pickable
		}
	}
	return nil
}

func (r *randImpl) BulkPick(drawCount int, pickables []Pickable) []Pickable {
	// ratioにマイナスが設定されると抽選ロジックがバグるので念のため除外する
	target := make([]Pickable, 0, len(pickables))
	for _, pickable := range pickables {
		if pickable.GetRatio() > 0 {
			target = append(target, pickable)
		}
	}
	if len(target) == 0 {
		return []Pickable{}
	}

	var totalWeight int64
	for _, pickable := range target {
		totalWeight += pickable.GetRatio()
	}
	if totalWeight <= 0 {
		return []Pickable{}
	}
	result := make([]Pickable, 0, drawCount)
	for i := 0; i < drawCount; i++ {
		random := r.NextInt63n(totalWeight)
		var temp int64
		for _, pickable := range target {
			temp += pickable.GetRatio()
			if temp > random {
				result = append(result, pickable)
				break
			}
		}
	}
	return result
}

func (r *randImpl) BulkPickNoDuplication(drawCount int, pickables []Pickable) []Pickable {
	// ratioにマイナスが設定されると抽選ロジックがバグるので念のため除外する
	target := make([]Pickable, 0, len(pickables))
	for _, pickable := range pickables {
		if pickable.GetRatio() > 0 {
			target = append(target, pickable)
		}
	}
	if len(target) == 0 {
		return []Pickable{}
	}

	result := make([]Pickable, 0, drawCount)
	for i := 0; i < drawCount; i++ {
		if len(target) == 0 {
			break
		}
		var totalWeight int64
		for _, pickable := range target {
			totalWeight += pickable.GetRatio()
		}
		if totalWeight <= 0 {
			break
		}
		random := r.NextInt63n(totalWeight)
		var temp int64
		var found bool
		remain := make([]Pickable, 0, len(target)-1)
		for _, pickable := range target {
			temp += pickable.GetRatio()
			if temp > random && !found {
				result = append(result, pickable)
				found = true
				continue
			}
			remain = append(remain, pickable)
		}
		target = remain
	}
	return result
}

func (r *randImpl) NewRandomUUID() (string, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return newUUID.String(), nil
}

func (r *randImpl) Shuffle(n int, swap func(i, j int)) {
	// ジェネリクスができたら直接スライス貰ってここでシャッフルしてしまいたいところ
	r.Rand.Shuffle(n, swap)
}
