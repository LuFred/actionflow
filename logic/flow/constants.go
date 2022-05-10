package flow

import (
	"strconv"
)

type ActionType int32

const (
	ActionTypeBlank ActionType = iota + 100
	ActionTypeHTTP
	ActionTypePage
)

var ActionTypeId = map[string]int32{
	"Blank": 100,
	"HTTP":  101,
	"Page":  102,
}
var ActionTypeName = map[int32]string{
	100: "Blank",
	101: "HTTP",
	102: "Page",
}

func (x ActionType) String() string {
	s, ok := ActionTypeName[int32(x)]
	if ok {
		return s
	}

	return strconv.Itoa(int(x))
}
