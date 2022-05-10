package flowdto

import "actionflow/dto"

type FlowJobDto struct {
	Id             string               `json:"id" description:""`
	FlowId         string               `json:"flowId" description:"流程id"`
	FlowVariable   string               `json:"flowVariable" description:"流程变量"`
	Timeout        int64                `json:"timeout" description:"超时时间"`
	Status         string               `json:"status" description:"任务状态"`
	Message        string               `json:"message" description:"结果消息"`
	RetryOptionsId *dto.RetryOptionsDto `json:"retryOptionsId" description:"重试策略"`
	CreatedBy      string               `json:"createdBy" description:"创建人id"`
	CreatedAt      int64                `json:"createdAt" description:"创建时间"`
	ModifiedBy     string               `json:"modifiedBy" description:"修改人id"`
	ModifiedAt     int64                `json:"modifiedAt" description:"修改时间"`
}
