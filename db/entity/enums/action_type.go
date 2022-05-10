package enums

import "strconv"

type ActionType int32

const (
	ACTION_TYPE_CONDITIONAL ActionType = 0

	ACTION_TYPE_ACTION ActionType = 1

	ACTION_TYPE_VARIABLE ActionType = 2
)

var ActionType_name = map[int32]string{
	0: "conditional",
	1: "action",
	2: "variable",
}

func (x ActionType) String() string {
	s, ok := ActionType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
