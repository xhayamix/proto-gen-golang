package time

import (
	"time"

	"github.com/xhayamix/proto-gen-golang/pkg/domain/enum"
)

type Resettable interface {
	ResetReadable
	GetResetVariable() string
	SetResetHour(hour int)
	SetResetMinute(minute int)
	SetResetWeek(week time.Weekday)
	SetResetDay(day int)
}

type ResetReadable interface {
	GetResetTimingType() enum.ResetTimingType
	GetResetHour() int
	GetResetMinute() int
	GetResetWeek() time.Weekday
	GetResetDay() int
}
