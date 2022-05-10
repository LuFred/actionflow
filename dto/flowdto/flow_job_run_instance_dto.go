package flowdto

import "actionflow/dto"

type FlowJobRunInstanceDto struct {
	Id             string               `json:"id" description:""`
	FlowId         string               `json:"flowId" description:"流程id"`
	JobId          string               `json:"jobId" description:"job id"`
	FlowVariable   string               `json:"flowVariable" description:"流程变量"`
	Timeout        int64                `json:"timeout" description:"超时时间"`
	Status         string               `json:"status" description:"状态 [pending | executing | success | failed | cancelled]"`
	Message        string               `json:"message" description:"结果消息"`
	RetryOptionsId *dto.RetryOptionsDto `json:"retryOptionsId" description:"重试策略"`
	ExecutedAt     int64                `json:"executedAt" description:"任务开始时间"`
	FinishedAt     int64                `json:"finishedAt" description:"任务结束时间"`
	CreatedBy      string               `json:"createdBy" description:"创建人id"`
	CreatedAt      int64                `json:"createdAt" description:"创建时间"`
	ModifiedBy     string               `json:"modifiedBy" description:"修改人id"`
	ModifiedAt     int64                `json:"modifiedAt" description:"修改时间"`
}
