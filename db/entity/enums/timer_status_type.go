package enums

import "strconv"

type TimerStatusType int32

const (
	TIMER_STATUS_TYPE_OPEN TimerStatusType = 0

	TIMER_STATUS_TYPE_CLOSED TimerStatusType = 1
)

var TimerStatusType_name = map[int32]string{
	0: "open",
	1: "closed",
}

func (x TimerStatusType) String() string {
	s, ok := TimerStatusType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
