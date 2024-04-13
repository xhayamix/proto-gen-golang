package time

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/xhayamix/proto-gen-golang/pkg/domain/constant"
	"github.com/xhayamix/proto-gen-golang/pkg/domain/enum"
)

var (
	DailyResetTimeRegExp   = regexp.MustCompile(`^(\d{1,2}):(\d{1,2})$`)
	WeeklyResetTimeRegExp  = regexp.MustCompile(`^(Sun|Mon|Tue|Wed|Thu|Fri|Sat) (\d{1,2}):(\d{1,2})$`)
	MonthlyResetTimeRegExp = regexp.MustCompile(`^(\d{1,2})(日|st|nd|rd|th) (\d{1,2}):(\d{1,2})$`)
	JST                    = time.FixedZone("Asia/Tokyo", 9*60*60)
)

const (
	dailyResetTimeRegExpResultSize   = 3
	weeklyResetTimeRegExpResultSize  = 4
	monthlyResetTimeRegExpResultSize = 5

	// secondsOfHour 1時間の秒数
	secondsOfHour = int64(time.Hour / time.Second)
	// secondsOfDay 1日の秒数
	secondsOfDay = int64(24 * time.Hour / time.Second)
)

// ToUnixMilli returns t as a Unix time, the number of milliseconds elapsed
// since January 1, 1970 UTC.
func ToUnixMilli(t *time.Time) int64 {
	if t == nil || t.IsZero() {
		return 0
	}
	return t.UnixMilli()
}

func UnixMilliToTime(unixMilli int64) time.Time {
	// 0のときはUTCの0ではなくtime.Time#IsZero()な値を返す
	if unixMilli == 0 {
		return time.Time{}
	}
	return time.UnixMilli(unixMilli)
}

func ToEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, time.Local)
}

func ToTime(ptr *time.Time) time.Time {
	if ptr == nil {
		return time.Time{}
	}
	return *ptr
}

func Format(t *time.Time, layout string) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.Format(layout)
}

// FormatBank BanKログ用に日付をフォーマットする
func FormatBank(t time.Time) string {
	return t.In(time.FixedZone("Asia/Tokyo", 9*60*60)).Format(time.DateTime)
}

// IsResetTime 最終更新時間/最終取得時間/etc..(lastTime)が現在時刻(now)と比べてリセットすべきものかどうか
func IsResetTime(lastTime, now time.Time, resettables ...ResetReadable) (bool, error) {
	_, prevResetTime, err := CalcPrevResetTime(now, resettables...)
	if err != nil {
		return false, err
	}
	// neverは何があってもfalse
	if prevResetTime.IsZero() {
		return false, nil
	}
	return lastTime.Before(prevResetTime), nil
}

// CalcPrevResetTime 前回リセット（可能な）日時を取得
func CalcPrevResetTime(now time.Time, resettables ...ResetReadable) (enum.ResetTimingType, time.Time, error) {
	var resetHour, resetMinute, resetDay int
	var resetWeek time.Weekday
	resetTimingType := enum.ResetTimingType_Never
	for _, r := range resettables {
		// 期間が一番短いリセットタイプを優先する. (期間の短さはenumのID昇順)
		if resetTimingType <= r.GetResetTimingType() {
			continue
		}
		resetTimingType = r.GetResetTimingType()
		resetHour = r.GetResetHour()
		resetMinute = r.GetResetMinute()
		resetWeek = r.GetResetWeek()
		resetDay = r.GetResetDay()
	}
	switch resetTimingType {
	case enum.ResetTimingType_Daily:
		todayResetTime := time.Date(now.Year(), now.Month(), now.Day(), resetHour, resetMinute, 0, 0, now.Location())
		if now.UnixNano() >= todayResetTime.UnixNano() {
			return resetTimingType, todayResetTime, nil
		}
		yesterdayResetTime := time.Date(now.Year(), now.Month(), now.Day()-1, resetHour, resetMinute, 0, 0, now.Location())
		return resetTimingType, yesterdayResetTime, nil
	case enum.ResetTimingType_Weekly:
		thisWeekResetTime := time.Date(now.Year(), now.Month(), now.Day(), resetHour, resetMinute, 0, 0, now.Location())
		thisWeekResetTime = thisWeekResetTime.Add(time.Duration(resetWeek-thisWeekResetTime.Weekday()) * 24 * time.Hour)
		if now.UnixNano() >= thisWeekResetTime.UnixNano() {
			return resetTimingType, thisWeekResetTime, nil
		}
		lastWeekResetTime := thisWeekResetTime.Add(-7 * 24 * time.Hour)
		return resetTimingType, lastWeekResetTime, nil
	case enum.ResetTimingType_Monthly:
		thisMonthResetTime := time.Date(now.Year(), now.Month(), resetDay, resetHour, resetMinute, 0, 0, now.Location())
		if now.UnixNano() >= thisMonthResetTime.UnixNano() {
			return resetTimingType, thisMonthResetTime, nil
		}
		lastMonthResetTime := time.Date(now.Year(), now.Month()-1, resetDay, resetHour, resetMinute, 0, 0, now.Location())
		return resetTimingType, lastMonthResetTime, nil
	case enum.ResetTimingType_Never:
		return resetTimingType, time.Time{}, nil
	default:
		return 0, time.Time{}, fmt.Errorf("unknown reset timing type: %v", resetTimingType)
	}
}

// CalcNextResetTime 次回リセット時間取得
func CalcNextResetTime(now time.Time, resettables ...ResetReadable) (enum.ResetTimingType, time.Time, error) {
	resetTimingType, prevResetTime, err := CalcPrevResetTime(now, resettables...)
	if err != nil {
		return 0, time.Time{}, err
	}
	switch resetTimingType {
	case enum.ResetTimingType_Daily:
		return resetTimingType, prevResetTime.AddDate(0, 0, 1), nil
	case enum.ResetTimingType_Weekly:
		return resetTimingType, prevResetTime.AddDate(0, 0, 7), nil
	case enum.ResetTimingType_Monthly:
		return resetTimingType, prevResetTime.AddDate(0, 1, 0), nil
	case enum.ResetTimingType_Never:
		return resetTimingType, time.Time{}, nil
	default:
		return 0, time.Time{}, fmt.Errorf("unknown reset timing type: %v", resetTimingType)
	}
}

func SetResetVariables(r Resettable) error {
	resetVariable := r.GetResetVariable()
	if resetVariable == "" {
		return nil
	}
	switch r.GetResetTimingType() {
	case enum.ResetTimingType_Daily:
		hour, minute, err := parseDailyResetTimeVariable(resetVariable)
		if err != nil {
			return err
		}
		r.SetResetHour(hour)
		r.SetResetMinute(minute)
	case enum.ResetTimingType_Weekly:
		week, hour, minute, err := parseWeeklyResetTimeVariable(resetVariable)
		if err != nil {
			return err
		}
		r.SetResetWeek(week)
		r.SetResetHour(hour)
		r.SetResetMinute(minute)
	case enum.ResetTimingType_Monthly:
		day, hour, minute, err := parseMonthlyResetTimeVariable(resetVariable)
		if err != nil {
			return err
		}
		r.SetResetDay(day)
		r.SetResetHour(hour)
		r.SetResetMinute(minute)
	case enum.ResetTimingType_Never:
		// 何もしない
		return nil
	}
	return nil
}

func parseDailyResetTimeVariable(resetTimingVariable string) (h, m int, e error) {
	times := DailyResetTimeRegExp.FindStringSubmatch(resetTimingVariable)
	if len(times) < dailyResetTimeRegExpResultSize {
		return 0, 0, fmt.Errorf("daily resetTimingVariable parse error. resetTimingVariable: %s", resetTimingVariable)
	}
	var err error
	// 正規表現で数値であることをチェックしているので、errはない想定
	hour, err := strconv.ParseInt(times[1], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("hour parse error. resetTimingVariable: %s, times[1]: %v", resetTimingVariable, times[1])
	}
	minute, err := strconv.ParseInt(times[2], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("minute parse error. resetTimingVariable: %s, times[2]: %v", resetTimingVariable, times[2])
	}
	return int(hour), int(minute), nil
}

func parseWeeklyResetTimeVariable(resetTimingVariable string) (w time.Weekday, h, m int, e error) {
	times := WeeklyResetTimeRegExp.FindStringSubmatch(resetTimingVariable)
	if len(times) < weeklyResetTimeRegExpResultSize {
		return time.Sunday, 0, 0, fmt.Errorf("weekly resetTimingVariable parse error. resetTimingVariable: %s", resetTimingVariable)
	}
	var weekday time.Weekday
	switch constant.Weekday(times[1]) {
	case constant.WeekdaySun:
		weekday = time.Sunday
	case constant.WeekdayMon:
		weekday = time.Monday
	case constant.WeekdayThe:
		weekday = time.Tuesday
	case constant.WeekdayWed:
		weekday = time.Wednesday
	case constant.WeekdayThu:
		weekday = time.Thursday
	case constant.WeekdayFri:
		weekday = time.Friday
	case constant.WeekdaySat:
		weekday = time.Saturday
	default:
		return time.Sunday, 0, 0, fmt.Errorf("unknown week: resetTimingVariable: %s, times[1]: %v", resetTimingVariable, times[1])
	}
	var err error
	// 正規表現で数値であることをチェックしているので、errはない想定
	hour, err := strconv.ParseInt(times[2], 10, 64)
	if err != nil {
		return time.Sunday, 0, 0, fmt.Errorf("hour parse error. resetTimingVariable: %s, times[1]: %v", resetTimingVariable, times[2])
	}
	minute, err := strconv.ParseInt(times[3], 10, 64)
	if err != nil {
		return time.Sunday, 0, 0, fmt.Errorf("minute parse error. resetTimingVariable: %s, times[2]: %v", resetTimingVariable, times[3])
	}
	return weekday, int(hour), int(minute), nil
}

func parseMonthlyResetTimeVariable(resetTimingVariable string) (d, h, m int, e error) {
	times := MonthlyResetTimeRegExp.FindStringSubmatch(resetTimingVariable)
	if len(times) < monthlyResetTimeRegExpResultSize {
		return 0, 0, 0, fmt.Errorf("monthly resetTimingVariable parse error. resetTimingVariable: %s", resetTimingVariable)
	}
	var err error
	// 正規表現で数値であることをチェックしているので、errはない想定
	day, err := strconv.ParseInt(times[1], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("day parse error. resetTimingVariable: %s, times[0]: %v", resetTimingVariable, times[1])
	}
	hour, err := strconv.ParseInt(times[3], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("hour parse error. resetTimingVariable: %s, times[1]: %v", resetTimingVariable, times[3])
	}
	minute, err := strconv.ParseInt(times[4], 10, 64)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("minute parse error. resetTimingVariable: %s, times[2]: %v", resetTimingVariable, times[4])
	}
	return int(day), int(hour), int(minute), nil
}

func TruncateHour(t time.Time, resetMinute int) time.Time {
	tt := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), resetMinute, 0, 0, t.Location())
	if tt.After(t) {
		tt = tt.Add(-time.Hour)
	}
	return tt
}

func CalcDiffHours(t1, t2 time.Time, resetMinute int) int32 {
	return int32(int64(TruncateHour(t2, resetMinute).Sub(TruncateHour(t1, resetMinute)).Seconds()) / secondsOfHour)
}

// IsSameHour 同じ時間かどうか
func IsSameHour(t1, t2 time.Time, resetMinute int) bool {
	return CalcDiffHours(t1, t2, resetMinute) == 0
}

// TruncateDay 渡した日付のゲーム内時刻での切り捨てた時間を返す
func TruncateDay(t time.Time, resetHour, resetMinute int) time.Time {
	tt := time.Date(t.Year(), t.Month(), t.Day(), resetHour, resetMinute, 0, 0, t.Location())
	if tt.After(t) {
		tt = tt.Add(-24 * time.Hour)
	}
	return tt
}

// CalcDiffDays 渡した日付のゲーム内時刻での日数差を返す
func CalcDiffDays(t1, t2 time.Time, resetHour, resetMinute int) int32 {
	return int32(int64(TruncateDay(t2, resetHour, resetMinute).Sub(TruncateDay(t1, resetHour, resetMinute)).Seconds()) / secondsOfDay)
}

// IsSameDay 同じ日付かどうか
func IsSameDay(t1, t2 time.Time, resetHour, resetMinute int) bool {
	return CalcDiffDays(t1, t2, resetHour, resetMinute) == 0
}

// TruncateMonth 渡した日付をゲーム内時刻で切り捨てた月を返す
func TruncateMonth(t time.Time, resetHour, resetMinute int) time.Time {
	tt := time.Date(t.Year(), t.Month(), 1, resetHour, resetMinute, 0, 0, t.Location())
	if tt.After(t) {
		tt = tt.AddDate(0, -1, 0)
	}
	return tt
}

// IsSameWeek 同じ週かどうか
func IsSameWeek(t1, t2 time.Time, resetHour, resetMinute int, resetWeekday time.Weekday) bool {
	tt1 := time.Date(t1.Year(), t1.Month(), t1.Day(), resetHour, resetMinute, 0, 0, t1.Location())
	tt1 = tt1.Add(time.Duration(resetWeekday-tt1.Weekday()) * 24 * time.Hour)
	if tt1.After(t1) {
		tt1 = tt1.Add(-7 * 24 * time.Hour)
	}
	tt2 := time.Date(t2.Year(), t2.Month(), t2.Day(), resetHour, resetMinute, 0, 0, t2.Location())
	tt2 = tt2.Add(time.Duration(resetWeekday-tt2.Weekday()) * 24 * time.Hour)
	if tt2.After(t2) {
		tt2 = tt2.Add(-7 * 24 * time.Hour)
	}
	return tt1.Equal(tt2)
}

// IsSameWeekday 入力日時が指定された曜日かどうか
func IsSameWeekday(t time.Time, weekdayStr string, resetHour, resetMinute int) bool {
	var weekday time.Weekday
	switch constant.Weekday(weekdayStr) {
	case constant.WeekdaySun:
		weekday = time.Sunday
	case constant.WeekdayMon:
		weekday = time.Monday
	case constant.WeekdayThe:
		weekday = time.Tuesday
	case constant.WeekdayWed:
		weekday = time.Wednesday
	case constant.WeekdayThu:
		weekday = time.Thursday
	case constant.WeekdayFri:
		weekday = time.Friday
	case constant.WeekdaySat:
		weekday = time.Saturday
	default:
		return false
	}

	tt := time.Date(t.Year(), t.Month(), t.Day(), resetHour, resetMinute, 0, 0, t.Location())
	if tt.After(t) {
		tt = tt.Add(-24 * time.Hour)
	}

	return tt.Weekday() == weekday
}

func IsInTerm(start, end, target time.Time) bool {
	if !start.IsZero() && target.Before(start) {
		return false
	}
	if !end.IsZero() && target.After(end) {
		return false
	}
	return true
}

func IsOver(end, target time.Time) bool {
	if end.IsZero() {
		return false
	}
	return end.UnixNano() < target.UnixNano()
}

func CalcAge(now, birthday time.Time) (int32, error) {
	if birthday.IsZero() {
		return 0, nil
	}
	if now.UnixNano() < birthday.UnixNano() {
		return 0, fmt.Errorf("誕生日が未来です. birthday=%s", birthday.Format(time.DateOnly))
	}
	age := int32(now.Year() - birthday.Year())
	if now.Month() < birthday.Month() || (now.Month() == birthday.Month() && now.Day() < birthday.Day()) {
		age--
	}
	return age, nil
}

func IsLeapYear(t time.Time) bool {
	// 閏年
	// 400で割り切れる年
	// 4で割り切れて、しかも100では割り切れない年
	return t.Year()%400 == 0 || (t.Year()%4 == 0 && t.Year()%100 != 0)
}
