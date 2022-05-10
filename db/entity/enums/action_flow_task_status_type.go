package enums

import "strconv"

type ActionFlowTaskStatusType int32

const (
	ACTION_FLOW_TASK_STATUS_TYPE_PENDING TimerStatusType = 0

	ACTION_FLOW_TASK_STATUS_TYPE_EXECUTING TimerStatusType = 1

	ACTION_FLOW_TASK_STATUS_TYPE_SUCCESS TimerStatusType = 2

	ACTION_FLOW_TASK_STATUS_TYPE_FAILED TimerStatusType = 3

	ACTION_FLOW_TASK_STATUS_TYPE_CANCELLED TimerStatusType = 4
)

var ActionFlowTaskStatusType_name = map[int32]string{
	0: "pending",
	1: "executing",
	2: "success",
	3: "failed",
	4: "cancelled",
}

func (x ActionFlowTaskStatusType) String() string {
	s, ok := ActionFlowTaskStatusType_name[int32(x)]
	if ok {
		return s
	}
	return strconv.Itoa(int(x))
}
