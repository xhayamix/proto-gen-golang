package time

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/xhayamix/proto-gen-golang/pkg/domain/enum"
)

type testResetReadable struct {
	resetTimingType enum.ResetTimingType
	resetHour       int
	resetMinute     int
	resetWeek       time.Weekday
	resetDay        int
}

func (r *testResetReadable) GetResetTimingType() enum.ResetTimingType { return r.resetTimingType }
func (r *testResetReadable) GetResetHour() int                        { return r.resetHour }
func (r *testResetReadable) GetResetMinute() int                      { return r.resetMinute }
func (r *testResetReadable) GetResetWeek() time.Weekday               { return r.resetWeek }
func (r *testResetReadable) GetResetDay() int                         { return r.resetDay }

func Test_IsResetTime(t *testing.T) {
	tests := map[string]struct {
		now      time.Time
		lastTime time.Time
		r        ResetReadable
		expected bool
		err      error
	}{
		"デイリー: リセット対象": {
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Daily,
				resetHour:       0,
				resetMinute:     0,
			},
			true,
			nil,
		},
		"デイリー: リセット対象: 同日": {
			time.Date(2021, 1, 2, 1, 0, 0, 0, time.Local),
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Daily,
				resetHour:       1,
				resetMinute:     0,
			},
			true,
			nil,
		},
		"デイリー: リセット対象外": {
			time.Date(2021, 1, 2, 2, 0, 0, 0, time.Local),
			time.Date(2021, 1, 2, 1, 0, 0, 0, time.Local),
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Daily,
				resetHour:       0,
				resetMinute:     0,
			},
			false,
			nil,
		},
		"デイリー: リセット対象外: 同タイミング": {
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Daily,
				resetHour:       0,
				resetMinute:     0,
			},
			false,
			nil,
		},
		"ウィークリー: リセット対象": {
			time.Date(2021, 1, 8, 0, 0, 0, 0, time.Local),
			time.Date(2021, 1, 4, 0, 0, 0, 0, time.Local),
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Weekly,
				resetWeek:       time.Wednesday,
				resetHour:       0,
				resetMinute:     0,
			},
			true,
			nil,
		},
		"ウィークリー: リセット対象: 1週間以上前のlastTime": {
			time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
			time.Date(2020, 12, 29, 0, 0, 0, 0, time.Local),
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Weekly,
				resetWeek:       time.Wednesday,
				resetHour:       0,
				resetMinute:     0,
			},
			true,
			nil,
		},
		"ウィークリー: リセット対象外: リセット曜日前": {
			time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
			time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Weekly,
				resetWeek:       time.Wednesday,
				resetHour:       0,
				resetMinute:     0,
			},
			false,
			nil,
		},
		"ウィークリー: リセット対象外: リセット曜日後": {
			time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
			time.Date(2021, 1, 6, 0, 0, 0, 0, time.Local),
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Weekly,
				resetWeek:       time.Wednesday,
				resetHour:       0,
				resetMinute:     0,
			},
			false,
			nil,
		},
		"ウィークリー: リセット対象外: ぎり1週間経ってないlastTime": {
			time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
			time.Date(2020, 12, 30, 0, 0, 0, 0, time.Local),
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Weekly,
				resetWeek:       time.Wednesday,
				resetHour:       0,
				resetMinute:     0,
			},
			false,
			nil,
		},
		"マンスリー: リセット対象": {
			time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
			time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Monthly,
				resetDay:        10,
				resetHour:       0,
				resetMinute:     0,
			},
			true,
			nil,
		},
		"マンスリー: リセット対象: 1ヶ月以上前のlastTime": {
			time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
			time.Date(2020, 1, 10, 0, 0, 0, 0, time.Local),
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Monthly,
				resetDay:        10,
				resetHour:       0,
				resetMinute:     0,
			},
			true,
			nil,
		},
		"マンスリー: リセット対象外: リセット日前": {
			time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
			time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Monthly,
				resetDay:        10,
				resetHour:       0,
				resetMinute:     0,
			},
			false,
			nil,
		},
		"マンスリー: リセット対象外: リセット日後": {
			time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
			time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Monthly,
				resetDay:        10,
				resetHour:       0,
				resetMinute:     0,
			},
			false,
			nil,
		},
		"マンスリー: リセット対象外: ぎり1ヶ月経ってないlastTime": {
			time.Date(2021, 1, 9, 0, 0, 0, 0, time.Local),
			time.Date(2020, 12, 10, 0, 0, 0, 0, time.Local),
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Monthly,
				resetDay:        10,
				resetHour:       0,
				resetMinute:     0,
			},
			false,
			nil,
		},
		"Never: 常にリセット対象外": {
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			time.Time{},
			&testResetReadable{
				resetTimingType: enum.ResetTimingType_Never,
			},
			false,
			nil,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			result, err := IsResetTime(tt.lastTime, tt.now, tt.r)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_CalcPrevResetTime(t *testing.T) {
	tests := map[string]struct {
		now                     time.Time
		rs                      []ResetReadable
		expected                time.Time
		expectedResetTimingType enum.ResetTimingType
		err                     error
	}{
		"デイリー: リセット時間前": {
			time.Date(2021, 1, 2, 0, 0, 0, 0, time.Local),
			[]ResetReadable{
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Daily,
					resetHour:       10,
					resetMinute:     30,
				},
			},
			time.Date(2021, 1, 1, 10, 30, 0, 0, time.Local),
			enum.ResetTimingType_Daily,
			nil,
		},
		"デイリー: リセット時間後": {
			time.Date(2021, 1, 2, 12, 0, 0, 0, time.Local),
			[]ResetReadable{
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Daily,
					resetHour:       10,
					resetMinute:     30,
				},
			},
			time.Date(2021, 1, 2, 10, 30, 0, 0, time.Local),
			enum.ResetTimingType_Daily,
			nil,
		},
		"ウィークリー: リセット曜日前": {
			time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
			[]ResetReadable{
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Weekly,
					resetWeek:       time.Wednesday,
					resetHour:       10,
					resetMinute:     30,
				},
			},
			time.Date(2020, 12, 30, 10, 30, 0, 0, time.Local),
			enum.ResetTimingType_Weekly,
			nil,
		},
		"ウィークリー: リセット曜日後": {
			time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
			[]ResetReadable{
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Weekly,
					resetWeek:       time.Wednesday,
					resetHour:       10,
					resetMinute:     30,
				},
			},
			time.Date(2021, 1, 6, 10, 30, 0, 0, time.Local),
			enum.ResetTimingType_Weekly,
			nil,
		},
		"マンスリー: リセット日前": {
			time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
			[]ResetReadable{
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Monthly,
					resetDay:        15,
					resetHour:       10,
					resetMinute:     30,
				},
			},
			time.Date(2020, 12, 15, 10, 30, 0, 0, time.Local),
			enum.ResetTimingType_Monthly,
			nil,
		},
		"マンスリー: リセット日後": {
			time.Date(2021, 1, 16, 0, 0, 0, 0, time.Local),
			[]ResetReadable{
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Monthly,
					resetDay:        15,
					resetHour:       10,
					resetMinute:     30,
				},
			},
			time.Date(2021, 1, 15, 10, 30, 0, 0, time.Local),
			enum.ResetTimingType_Monthly,
			nil,
		},
		"2種類のResetReadableが入ったとき、スパンが短いほうを使用する": {
			time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
			[]ResetReadable{
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Monthly,
					resetDay:        1,
					resetHour:       0,
					resetMinute:     0,
				},
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Weekly,
					resetWeek:       time.Wednesday,
					resetHour:       10,
					resetMinute:     30,
				},
			},
			time.Date(2020, 12, 30, 10, 30, 0, 0, time.Local),
			enum.ResetTimingType_Weekly,
			nil,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			resetTimingType, result, err := CalcPrevResetTime(tt.now, tt.rs...)
			assert.Equal(t, tt.expectedResetTimingType, resetTimingType)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.err, err)
		})
	}
}

func Test_CalcNextResetTime(t *testing.T) {
	tests := map[string]struct {
		now                     time.Time
		rs                      []ResetReadable
		expected                time.Time
		expectedResetTimingType enum.ResetTimingType
		err                     error
	}{
		"デイリー: リセット時間前": {
			time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local),
			[]ResetReadable{
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Daily,
					resetHour:       10,
					resetMinute:     30,
				},
			},
			time.Date(2021, 1, 1, 10, 30, 0, 0, time.Local),
			enum.ResetTimingType_Daily,
			nil,
		},
		"デイリー: リセット時間後": {
			time.Date(2021, 1, 1, 12, 0, 0, 0, time.Local),
			[]ResetReadable{
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Daily,
					resetHour:       10,
					resetMinute:     30,
				},
			},
			time.Date(2021, 1, 2, 10, 30, 0, 0, time.Local),
			enum.ResetTimingType_Daily,
			nil,
		},
		"ウィークリー: リセット曜日前": {
			time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
			[]ResetReadable{
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Weekly,
					resetWeek:       time.Wednesday,
					resetHour:       10,
					resetMinute:     30,
				},
			},
			time.Date(2021, 1, 6, 10, 30, 0, 0, time.Local),
			enum.ResetTimingType_Weekly,
			nil,
		},
		"ウィークリー: リセット曜日後": {
			time.Date(2021, 1, 7, 0, 0, 0, 0, time.Local),
			[]ResetReadable{
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Weekly,
					resetWeek:       time.Wednesday,
					resetHour:       10,
					resetMinute:     30,
				},
			},
			time.Date(2021, 1, 13, 10, 30, 0, 0, time.Local),
			enum.ResetTimingType_Weekly,
			nil,
		},
		"マンスリー: リセット日前": {
			time.Date(2021, 1, 10, 0, 0, 0, 0, time.Local),
			[]ResetReadable{
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Monthly,
					resetDay:        15,
					resetHour:       10,
					resetMinute:     30,
				},
			},
			time.Date(2021, 1, 15, 10, 30, 0, 0, time.Local),
			enum.ResetTimingType_Monthly,
			nil,
		},
		"マンスリー: リセット日後": {
			time.Date(2021, 1, 16, 0, 0, 0, 0, time.Local),
			[]ResetReadable{
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Monthly,
					resetDay:        15,
					resetHour:       10,
					resetMinute:     30,
				},
			},
			time.Date(2021, 2, 15, 10, 30, 0, 0, time.Local),
			enum.ResetTimingType_Monthly,
			nil,
		},
		"2種類のResetReadableが入ったとき、スパンが短いほうを使用する": {
			time.Date(2021, 1, 5, 0, 0, 0, 0, time.Local),
			[]ResetReadable{
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Monthly,
					resetDay:        1,
					resetHour:       0,
					resetMinute:     0,
				},
				&testResetReadable{
					resetTimingType: enum.ResetTimingType_Weekly,
					resetWeek:       time.Wednesday,
					resetHour:       10,
					resetMinute:     30,
				},
			},
			time.Date(2021, 1, 6, 10, 30, 0, 0, time.Local),
			enum.ResetTimingType_Weekly,
			nil,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			resetTimingType, result, err := CalcNextResetTime(tt.now, tt.rs...)
			assert.Equal(t, tt.expectedResetTimingType, resetTimingType)
			assert.Equal(t, tt.expected, result)
			assert.Equal(t, tt.err, err)
		})
	}
}
func Test_TruncateHour(t *testing.T) {
	for name, tt := range map[string]struct {
		t1           time.Time
		resetMinutes int
		expected     time.Time
	}{
		"10:00に設定": {
			t1:           time.Date(2022, 1, 1, 10, 40, 10, 100, time.Local),
			resetMinutes: 0,
			expected:     time.Date(2022, 1, 1, 10, 0, 0, 0, time.Local),
		},
		"10:30に設定": {
			t1:           time.Date(2022, 1, 1, 10, 40, 10, 100, time.Local),
			resetMinutes: 30,
			expected:     time.Date(2022, 1, 1, 10, 30, 0, 0, time.Local),
		},
		"10:30に設定(リセットタイミング前なので1時間早める)": {
			t1:           time.Date(2022, 1, 1, 10, 10, 10, 100, time.Local),
			resetMinutes: 30,
			expected:     time.Date(2022, 1, 1, 9, 30, 0, 0, time.Local),
		},
	} {
		tt := tt
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, TruncateHour(tt.t1, tt.resetMinutes))
		})
	}
}

func Test_CalcDiffHours(t *testing.T) {
	for name, tt := range map[string]struct {
		t1           time.Time
		t2           time.Time
		resetMinutes int
		expected     int32
	}{
		"同日、同時間、同分 resetMinutesなし": {
			t1:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			resetMinutes: 0,
			expected:     0,
		},
		"同日、同時間、t1が前の分 resetMinutesなし": {
			t1:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 1, 59, 0, 0, time.Local),
			resetMinutes: 0,
			expected:     0,
		},
		"同日、同時間、t1が後の分 resetMinutesなし": {
			t1:           time.Date(2022, 1, 1, 1, 59, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			resetMinutes: 0,
			expected:     0,
		},
		"同日、t1が前の時間、同分 resetMinutesなし": {
			t1:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 2, 0, 0, 0, time.Local),
			resetMinutes: 0,
			expected:     1,
		},
		"同日、t1が後の時間、同分 resetMinutesなし": {
			t1:           time.Date(2022, 1, 1, 2, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			resetMinutes: 0,
			expected:     -1,
		},
		"t1が前の日、同時間、同分 resetMinutesなし": {
			t1:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 2, 1, 0, 0, 0, time.Local),
			resetMinutes: 0,
			expected:     24,
		},
		"t1が後の日、同時間、同分 resetMinutesなし": {
			t1:           time.Date(2022, 1, 2, 1, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			resetMinutes: 0,
			expected:     -24,
		},

		"同日、同時間、同分 resetMinutesあり": {
			t1:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			resetMinutes: 30,
			expected:     0,
		},

		"同日、同時間、t1が前の分 resetMinutesあり": {
			t1:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 1, 29, 0, 0, time.Local),
			resetMinutes: 30,
			expected:     0,
		},
		"同日、同時間、t1が後の分 resetMinutesあり": {
			t1:           time.Date(2022, 1, 1, 1, 29, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			resetMinutes: 30,
			expected:     0,
		},
		"同日、同時間、t1がリセットタイミングを跨いで前の分 resetMinutesあり": {
			t1:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 1, 59, 0, 0, time.Local),
			resetMinutes: 30,
			expected:     1,
		},
		"同日、同時間、t1がリセットタイミングを跨いで後の分 resetMinutesあり": {
			t1:           time.Date(2022, 1, 1, 1, 59, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			resetMinutes: 30,
			expected:     -1,
		},
		"同日、t1が前の時間、同分 resetMinutesあり": {
			t1:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 2, 0, 0, 0, time.Local),
			resetMinutes: 30,
			expected:     1,
		},
		"同日、t1が後の時間、同分 resetMinutesあり": {
			t1:           time.Date(2022, 1, 1, 2, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			resetMinutes: 30,
			expected:     -1,
		},
		"t1が前の日、同時間、同分 resetMinutesあり": {
			t1:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 2, 1, 0, 0, 0, time.Local),
			resetMinutes: 30,
			expected:     24,
		},
		"t1が後の日、同時間、同分 resetMinutesあり": {
			t1:           time.Date(2022, 1, 2, 1, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			resetMinutes: 30,
			expected:     -24,
		},
	} {
		tt := tt
		t.Run(name, func(t *testing.T) {
			days := CalcDiffHours(tt.t1, tt.t2, tt.resetMinutes)
			assert.Equal(t, tt.expected, days)
		})
	}
}

func Test_IsSameHour(t *testing.T) {
	tests := map[string]struct {
		t1       time.Time
		t2       time.Time
		expected bool
	}{
		"時間差が0ならtrue": {
			t1:       time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			t2:       time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			expected: true,
		},
		"時間差が0でないならfalse": {
			t1:       time.Date(2022, 1, 1, 1, 0, 0, 0, time.Local),
			t2:       time.Date(2022, 1, 1, 2, 0, 0, 0, time.Local),
			expected: false,
		},
	}
	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			result := IsSameHour(tt.t1, tt.t2, 0)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_TruncateDay(t *testing.T) {
	for name, tt := range map[string]struct {
		t1           time.Time
		resetHour    int
		resetMinutes int
		expected     time.Time
	}{
		"0:00 リセット": {
			t1:           time.Date(2022, 1, 1, 10, 0, 0, 0, time.Local),
			resetHour:    0,
			resetMinutes: 0,
			expected:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
		},
		"4:30 リセット": {
			t1:           time.Date(2022, 1, 1, 15, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 30,
			expected:     time.Date(2022, 1, 1, 4, 30, 0, 0, time.Local),
		},
		"4:30 リセット_リセットタイミングを過ぎていないため日付が前日になっている": {
			t1:           time.Date(2022, 1, 2, 0, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 30,
			expected:     time.Date(2022, 1, 1, 4, 30, 0, 0, time.Local),
		},
	} {
		tt := tt
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, TruncateDay(tt.t1, tt.resetHour, tt.resetMinutes))
		})
	}
}

func Test_TruncateMonth(t *testing.T) {
	for name, tt := range map[string]struct {
		t1           time.Time
		resetHour    int
		resetMinutes int
		expected     time.Time
	}{
		"0:00 リセット": {
			t1:           time.Date(2022, 1, 2, 10, 0, 0, 0, time.Local),
			resetHour:    0,
			resetMinutes: 0,
			expected:     time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
		},
		"4:30 リセット": {
			t1:           time.Date(2022, 1, 2, 15, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 30,
			expected:     time.Date(2022, 1, 1, 4, 30, 0, 0, time.Local),
		},
		"4:30 リセット_リセットタイミングを過ぎていないため日付が前月になっている": {
			t1:           time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 30,
			expected:     time.Date(2021, 12, 1, 4, 30, 0, 0, time.Local),
		},
	} {
		tt := tt
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.expected, TruncateMonth(tt.t1, tt.resetHour, tt.resetMinutes))
		})
	}
}

func Test_CalcDiffDays(t *testing.T) {
	for name, tt := range map[string]struct {
		t1           time.Time
		t2           time.Time
		resetHour    int
		resetMinutes int
		expected     int32
	}{
		"t1.Date < t2.Date, t1.Hour > ResetHour, t2.Hour > ResetHour, t1.Hour > t2.Hour": {
			t1:           time.Date(2022, 1, 1, 15, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 5, 10, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     4,
		},
		"t1.Date < t2.Date, t1.Hour > ResetHour, t2.Hour = ResetHour, t1.Hour > t2.Hour": {
			t1:           time.Date(2022, 1, 1, 15, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 5, 4, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     4,
		},
		"t1.Date < t2.Date, t1.Hour > ResetHour, t2.Hour > ResetHour, t1.Hour < t2.Hour": {
			t1:           time.Date(2022, 1, 1, 10, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 5, 15, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     4,
		},
		"t1.Date < t2.Date, t1.Hour = ResetHour, t2.Hour > ResetHour, t1.Hour < t2.Hour": {
			t1:           time.Date(2022, 1, 1, 4, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 5, 15, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     4,
		},
		"t1.Date < t2.Date, t1.Hour < ResetHour, t2.Hour < ResetHour, t1.Hour > t2.Hour": {
			t1:           time.Date(2022, 1, 1, 3, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 5, 2, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     4,
		},
		"t1.Date < t2.Date, t1.Hour < ResetHour, t2.Hour < ResetHour, t1.Hour < t2.Hour": {
			t1:           time.Date(2022, 1, 1, 2, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 5, 3, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     4,
		},
		"t1.Date < t2.Date, t1.Hour > ResetHour, t2.Hour < ResetHour": {
			t1:           time.Date(2022, 1, 1, 15, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 5, 3, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     3,
		},
		"t1.Date < t2.Date, t1.Hour < ResetHour, t2.Hour > ResetHour": {
			t1:           time.Date(2022, 1, 1, 3, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 5, 15, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     5,
		},
		"t1.Date = t2.Date, t1.Hour > ResetHour, t2.Hour > ResetHour, t1.Hour > t2.Hour": {
			t1:           time.Date(2022, 1, 1, 15, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 10, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     0,
		},
		"t1.Date = t2.Date, t1.Hour > ResetHour, t2.Hour > ResetHour, t1.Hour < t2.Hour": {
			t1:           time.Date(2022, 1, 1, 10, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 15, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     0,
		},
		"t1.Date = t2.Date, t1.Hour < ResetHour, t2.Hour < ResetHour, t1.Hour > t2.Hour": {
			t1:           time.Date(2022, 1, 1, 3, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 2, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     0,
		},
		"t1.Date = t2.Date, t1.Hour < ResetHour, t2.Hour < ResetHour, t1.Hour < t2.Hour": {
			t1:           time.Date(2022, 1, 1, 2, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 3, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     0,
		},
		"t1.Date = t2.Date, t1.Hour > ResetHour, t2.Hour < ResetHour": {
			t1:           time.Date(2022, 1, 1, 15, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 3, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     -1,
		},
		"t1.Date = t2.Date, t1.Hour < ResetHour, t2.Hour > ResetHour": {
			t1:           time.Date(2022, 1, 1, 3, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 15, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     1,
		},
		"t1.Date > t2.Date, t1.Hour > ResetHour, t2.Hour > ResetHour, t1.Hour > t2.Hour": {
			t1:           time.Date(2022, 1, 5, 15, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 10, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     -4,
		},
		"t1.Date > t2.Date, t1.Hour > ResetHour, t2.Hour = ResetHour, t1.Hour > t2.Hour": {
			t1:           time.Date(2022, 1, 5, 15, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 4, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     -4,
		},
		"t1.Date > t2.Date, t1.Hour > ResetHour, t2.Hour > ResetHour, t1.Hour < t2.Hour": {
			t1:           time.Date(2022, 1, 5, 10, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 15, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     -4,
		},
		"t1.Date > t2.Date, t1.Hour = ResetHour, t2.Hour > ResetHour, t1.Hour < t2.Hour": {
			t1:           time.Date(2022, 1, 5, 4, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 15, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     -4,
		},
		"t1.Date > t2.Date, t1.Hour < ResetHour, t2.Hour < ResetHour, t1.Hour > t2.Hour": {
			t1:           time.Date(2022, 1, 5, 3, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 2, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     -4,
		},
		"t1.Date > t2.Date, t1.Hour < ResetHour, t2.Hour < ResetHour, t1.Hour < t2.Hour": {
			t1:           time.Date(2022, 1, 5, 2, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 3, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     -4,
		},
		"t1.Date > t2.Date, t1.Hour < ResetHour, t2.Hour > ResetHour": {
			t1:           time.Date(2022, 1, 5, 3, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 10, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     -3,
		},
		"t1.Date > t2.Date, t1.Hour > ResetHour, t2.Hour < ResetHour": {
			t1:           time.Date(2022, 1, 5, 10, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 3, 0, 0, 0, time.Local),
			resetHour:    4,
			resetMinutes: 0,
			expected:     -5,
		},
		"ResetHour = 0, t1.Date > t2.Date": {
			t1:           time.Date(2022, 1, 5, 10, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 3, 0, 0, 0, time.Local),
			resetHour:    0,
			resetMinutes: 0,
			expected:     -4,
		},
		"ResetHour = 0, t1.Date < t2.Date": {
			t1:           time.Date(2022, 1, 1, 10, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 5, 3, 0, 0, 0, time.Local),
			resetHour:    0,
			resetMinutes: 0,
			expected:     4,
		},
		"ResetHour = 0, t1.Date = t2.Date": {
			t1:           time.Date(2022, 1, 1, 10, 0, 0, 0, time.Local),
			t2:           time.Date(2022, 1, 1, 3, 0, 0, 0, time.Local),
			resetHour:    0,
			resetMinutes: 0,
			expected:     0,
		},
	} {
		tt := tt
		t.Run(name, func(t *testing.T) {
			days := CalcDiffDays(tt.t1, tt.t2, tt.resetHour, tt.resetMinutes)
			assert.Equal(t, tt.expected, days)
		})
	}
}

func Test_IsSameDay(t *testing.T) {
	tests := map[string]struct {
		t1       time.Time
		t2       time.Time
		expected bool
	}{
		"日数差が0ならtrue": {
			t1:       time.Date(2022, 3, 1, 12, 0, 0, 0, time.Local),
			t2:       time.Date(2022, 3, 1, 12, 0, 0, 0, time.Local),
			expected: true,
		},
		"日数差が0でないならfalse": {
			t1:       time.Date(2022, 3, 1, 12, 0, 0, 0, time.Local),
			t2:       time.Date(2022, 3, 3, 12, 0, 0, 0, time.Local),
			expected: false,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			result := IsSameDay(tt.t1, tt.t2, 0, 0)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsSameWeek(t *testing.T) {
	for name, tt := range map[string]struct {
		t1       time.Time
		t2       time.Time
		expected bool
	}{
		"同じ週: どちらも日付変更線を超えている": {
			t1:       time.Date(2022, 3, 1, 6, 0, 0, 0, time.Local),
			t2:       time.Date(2022, 3, 6, 12, 0, 0, 0, time.Local),
			expected: true,
		},
		"同じ週: t1だけ日付変更線を超えている": {
			t1:       time.Date(2022, 2, 28, 6, 0, 0, 0, time.Local),
			t2:       time.Date(2022, 3, 7, 3, 0, 0, 0, time.Local),
			expected: true,
		},
		"同じ週: t2だけ日付変更線を超えている": {
			t1:       time.Date(2022, 3, 1, 3, 0, 0, 0, time.Local),
			t2:       time.Date(2022, 2, 28, 12, 0, 0, 0, time.Local),
			expected: true,
		},
		"同じ週: どちらも日付変更線を超えていない": {
			t1:       time.Date(2022, 3, 1, 3, 0, 0, 0, time.Local),
			t2:       time.Date(2022, 3, 7, 2, 0, 0, 0, time.Local),
			expected: true,
		},
		"別の週: どちらも日付変更線を超えている": {
			t1:       time.Date(2022, 2, 28, 6, 0, 0, 0, time.Local),
			t2:       time.Date(2022, 3, 7, 12, 0, 0, 0, time.Local),
			expected: false,
		},
		"別の週: t1だけ日付変更線を超えている": {
			t1:       time.Date(2022, 2, 28, 6, 0, 0, 0, time.Local),
			t2:       time.Date(2022, 3, 8, 3, 0, 0, 0, time.Local),
			expected: false,
		},
		"別の週: t2だけ日付変更線を超えている": {
			t1:       time.Date(2022, 3, 1, 3, 0, 0, 0, time.Local),
			t2:       time.Date(2022, 3, 7, 12, 0, 0, 0, time.Local),
			expected: false,
		},
		"別の週: どちらも日付変更線を超えていない": {
			t1:       time.Date(2022, 3, 1, 3, 0, 0, 0, time.Local),
			t2:       time.Date(2022, 3, 8, 4, 0, 0, 0, time.Local),
			expected: false,
		},
	} {
		tt := tt
		t.Run(name, func(t *testing.T) {
			result := IsSameWeek(tt.t1, tt.t2, 5, 0, time.Monday)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsSameWeekday(t *testing.T) {
	for name, tt := range map[string]struct {
		t1         time.Time
		weekdayStr string
		expected   bool
	}{
		"同じ週: 日付変更線を超えている": {
			t1:         time.Date(2022, 3, 1, 6, 0, 0, 0, time.Local),
			weekdayStr: "Tue",
			expected:   true,
		},
		"同じ週: 日付変更線を超えていない": {
			t1:         time.Date(2022, 3, 1, 4, 0, 0, 0, time.Local),
			weekdayStr: "Wed",
			expected:   false,
		},
	} {
		tt := tt
		t.Run(name, func(t *testing.T) {
			result := IsSameWeekday(tt.t1, tt.weekdayStr, 5, 0)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_IsInTerm(t *testing.T) {
	tests := map[string]struct {
		now      time.Time
		start    time.Time
		end      time.Time
		expected bool
	}{
		"期間前": {
			time.Date(2021, 12, 31, 23, 59, 59, 999, time.Local),
			time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
			time.Date(2022, 1, 31, 23, 59, 59, 0, time.Local),
			false,
		},
		"期間中": {
			time.Date(2022, 1, 1, 23, 59, 59, 999, time.Local),
			time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
			time.Date(2022, 1, 31, 23, 59, 59, 0, time.Local),
			true,
		},
		"期間後": {
			time.Date(2022, 2, 1, 0, 0, 0, 0, time.Local),
			time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
			time.Date(2022, 1, 31, 23, 59, 59, 0, time.Local),
			false,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			result := IsInTerm(tt.start, tt.end, tt.now)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_IsOver(t *testing.T) {
	tests := map[string]struct {
		now      time.Time
		end      time.Time
		expected bool
	}{
		"期間前": {
			time.Date(2021, 12, 31, 23, 59, 59, 999, time.Local),
			time.Date(2022, 1, 1, 0, 0, 0, 0, time.Local),
			false,
		},
		"期間中": {
			time.Date(2021, 12, 31, 23, 59, 59, 999, time.Local),
			time.Date(2021, 12, 31, 23, 59, 59, 999, time.Local),
			false,
		},
		"期間後": {
			time.Date(2021, 12, 31, 23, 59, 59, 999, time.Local),
			time.Date(2021, 12, 31, 23, 59, 59, 998, time.Local),
			true,
		},
	}

	for name, tt := range tests {
		tt := tt
		t.Run(name, func(t *testing.T) {
			result := IsOver(tt.end, tt.now)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCalcAge(t *testing.T) {
	for name, tt := range map[string]struct {
		target   time.Time
		birthday time.Time
		age      int32
		err      error
	}{
		"2000/7/01 ~ 2020/6/15 = 19": {
			target:   time.Date(2020, 6, 15, 0, 0, 0, 0, time.Local),
			birthday: time.Date(2000, 7, 1, 0, 0, 0, 0, time.Local),
			age:      19,
		},
		"2000/6/16 ~ 2020/6/15 = 20": {
			target:   time.Date(2020, 6, 15, 0, 0, 0, 0, time.Local),
			birthday: time.Date(2000, 6, 16, 0, 0, 0, 0, time.Local),
			age:      19,
		},
		"2000/6/15 ~ 2020/6/15 = 20": {
			target:   time.Date(2020, 6, 15, 0, 0, 0, 0, time.Local),
			birthday: time.Date(2000, 6, 15, 0, 0, 0, 0, time.Local),
			age:      20,
		},
		"2000/2/29 ~ 2020/6/15 = 20": {
			target:   time.Date(2020, 6, 15, 0, 0, 0, 0, time.Local),
			birthday: time.Date(2000, 2, 29, 0, 0, 0, 0, time.Local),
			age:      20,
		},
		"誕生日が閏年:2000/2/29 ~ 2020/2/28 = 19": {
			target:   time.Date(2020, 2, 28, 0, 0, 0, 0, time.Local),
			birthday: time.Date(2000, 2, 29, 0, 0, 0, 0, time.Local),
			age:      19,
		},
		"誕生日が閏年:2000/2/29 ~ 2020/2/29 = 20": {
			target:   time.Date(2020, 2, 29, 0, 0, 0, 0, time.Local),
			birthday: time.Date(2000, 2, 29, 0, 0, 0, 0, time.Local),
			age:      20,
		},
		"誕生日が閏年:2000/2/29 ~ 2021/2/28 = 20": {
			target:   time.Date(2021, 2, 28, 0, 0, 0, 0, time.Local),
			birthday: time.Date(2000, 2, 29, 0, 0, 0, 0, time.Local),
			age:      20,
		},
		"誕生日が閏年:2000/2/29 ~ 2021/3/1 = 21": {
			target:   time.Date(2021, 3, 1, 0, 0, 0, 0, time.Local),
			birthday: time.Date(2000, 2, 29, 0, 0, 0, 0, time.Local),
			age:      21,
		},
		"未来": {
			target:   time.Date(2020, 6, 15, 0, 0, 0, 0, time.Local),
			birthday: time.Date(2020, 6, 16, 0, 0, 0, 0, time.Local),
			err:      errors.New("誕生日が未来です. birthday=2020-06-16"),
		},
	} {
		tt := tt
		t.Run(name, func(t *testing.T) {
			age, err := CalcAge(tt.target, tt.birthday)
			if tt.err == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.err, err)
			}
			assert.Equal(t, tt.age, age)
		})
	}
}

func TestIsLeapYear(t *testing.T) {
	for name, tt := range map[string]struct {
		t   time.Time
		ret bool
	}{
		"閏年（400で割り切れる）": {
			t:   time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local),
			ret: true,
		},
		"閏年（4で割り切れるが100で割り切れない）": {
			t:   time.Date(2004, 1, 1, 0, 0, 0, 0, time.Local),
			ret: true,
		},
		"閏年じゃない（4でも100でも割り切れる）": {
			t:   time.Date(2100, 1, 1, 0, 0, 0, 0, time.Local),
			ret: false,
		},
	} {
		tt := tt
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.ret, IsLeapYear(tt.t))
		})
	}
}
