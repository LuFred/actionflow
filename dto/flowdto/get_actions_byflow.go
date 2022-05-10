package flowdto

type GetActionsByFlowIdRequest struct {
	FlowId string `json:"flowId" transport:"path" required:"true"  description:"流程Id"`
}

type GetActionsByFlowIdResponse struct {
	Data []*ActionDto `json:"data"`
}
