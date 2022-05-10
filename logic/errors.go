package logic

import (
	"fmt"
	"strconv"
)

type HttpError int32

const (
	HttpErrorNameAlreadyExists HttpError = 1000
	HttpErrorFlowNotFound      HttpError = 1001
	HttpErrorParameter         HttpError = 1002
	HttpErrorInvalidPreId      HttpError = 1003
	HttpErrorInvalidNextId     HttpError = 1004
	HttpErrorRepeatEdge        HttpError = 1005
	HttpErrorGraphHaveCycle    HttpError = 1006
	HttpErrorFlowJobNotFound   HttpError = 1007
)

var HttpErrorName = map[int32]string{
	1000: "Name already exists",
	1001: "flow not found",
	1002: "parameter error",
	1003: "invalid preId",
	1004: "invalid nextId",
	1005: "repeated edge",
	1006: "the graph have a cycle",
	1007: "flow job not found",
}

func (x HttpError) String() string {
	s, ok := HttpErrorName[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}

func (x HttpError) Error() string {
	return fmt.Sprintf("{\n    \"code\": %d,\n    \"message\": \"%s\"\n}", x, HttpErrorName[int32(x)])
}
