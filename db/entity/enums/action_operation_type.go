package enums

import "strconv"

type ActionHandleType int32

const (
	ACTION_HANDLE_TYPE_REST ActionHandleType = 0
)

var ActionHandleType_name = map[int32]string{
	0: "rest",
}

func (x ActionHandleType) String() string {
	s, ok := ActionHandleType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
