package flowdto

type CreateActionRequest struct {
	PreIds      []string `json:"preIds" transport:"body" required:"true" description:"父级id列表节点"`
	NextIds     []string `json:"nextIds" transport:"body" description:"子节点id列表"`
	Name        string   `json:"name" transport:"body" required:"true" description:"name"`
	DisplayName string   `json:"displayName" transport:"body" required:"true" description:"显示名称"`
	FlowId      string   `json:"flowId" transport:"body" required:"true"  description:"流程Id"`
	Type        string   `json:"type" transport:"body" required:"true"  description:"节点类型"`
	Command     string   `json:"command" transport:"body" required:"true" description:"操作命令"`
}

type CreateActionResponse struct {
	CreateFlowRequest
}
