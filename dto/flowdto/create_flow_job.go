package flowdto

type CreateFlowJobRequest struct {
	FlowId       string `json:"flowId" transport:"path" required:"true" description:"flow id"`
	FlowVariable string `json:"flowVariable" transport:"body" required:"false"  description:"流程变量"`
	Timeout      int64  `json:"timeout" transport:"body" required:"true" description:"超时时间"`
}
