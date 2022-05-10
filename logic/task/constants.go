package task

import "strconv"

type JobStatusType int32

const (
	JobStatusTypeActive JobStatusType = iota + 100
	JobStatusTypeDeleted
	JobStatusTypePause
)

type JobRunInstanceStatusType int32

const (
	JobRunInstanceStatusTypePending JobRunInstanceStatusType = iota + 100
	JobRunInstanceStatusTypePause
	JobRunInstanceStatusTypeExecuting
	JobRunInstanceStatusTypeSucceed
	JobRunInstanceStatusTypeFailed
	JobRunInstanceStatusTypeCancelled
)

var JobStatusTypeId = map[string]int32{
	"active":  100,
	"deleted": 101,
	"pause":   102,
}

var JobStatusTypeName = map[int32]string{
	100: "active",
	101: "deleted",
	102: "pause",
}

func (x JobStatusType) String() string {
	s, ok := JobStatusTypeName[int32(x)]
	if ok {
		return s
	}

	return strconv.Itoa(int(x))
}

var JobRunInstanceStatusTypeTypeId = map[string]int32{
	"pending":   100,
	"pause":     101,
	"executing": 102,
	"succeed":   103,
	"failed":    104,
	"cancelled": 105,
}

var JobRunInstanceStatusTypeName = map[int32]string{
	100: "pending",
	101: "pause",
	102: "executing",
	103: "succeed",
	104: "failed",
	105: "cancelled",
}

func (x JobRunInstanceStatusType) String() string {
	s, ok := JobRunInstanceStatusTypeName[int32(x)]
	if ok {
		return s
	}

	return strconv.Itoa(int(x))
}
