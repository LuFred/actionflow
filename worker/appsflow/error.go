package appsflow

import (
	"fmt"
	"strconv"
)

type ErrorType int32

const (
	ERROR_TYPE_TIMEOUT  ErrorType = 100
	ERROR_TYPE_CANCELED ErrorType = 101
	ERROR_TYPE_ACTIVITY ErrorType = 102
)

var ErrorType_name = map[int32]string{
	100: "timeout",
	101: "canceled",
	102: "activity",
}

var ErrorType_value = map[string]int32{
	"timeout":  100,
	"canceled": 101,
	"activity": 102,
}

func (x ErrorType) String() string {
	s, ok := ErrorType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}

type WorkflowError struct {
	msg     string
	errType ErrorType
}

func NewWorkflowError(msg string, errType ErrorType) error {
	applicationErr := &WorkflowError{
		msg:     msg,
		errType: errType}

	return applicationErr
}

func (e *WorkflowError) Error() string {
	msg := fmt.Sprintf("%s (type: %s)", e.msg, e.errType)
	return msg
}

func (e *WorkflowError) Type() string {
	return e.errType.String()
}
