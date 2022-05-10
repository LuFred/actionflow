package flowdto

import "actionflow/dto"

type ActionDto struct {
	PreIds         []string             `json:"preIds" description:"父节点id"`
	Id             string               `json:"id" description:""`
	FlowId         string               `json:"flowId" description:"流程id"`
	Name           string               `json:"name" description:"name"`
	DisplayName    string               `json:"displayName" description:"显示名称"`
	Type           string               `json:"type" description:"节点类型"`
	Command        string               `json:"command" description:"操作命令"`
	RetryOptionsId *dto.RetryOptionsDto `json:"retryOptionsId" description:"重试策略"`
	CreatedBy      string               `json:"createdBy" description:"创建人id"`
	CreatedAt      int64                `json:"createdAt" description:"创建时间"`
	ModifiedBy     string               `json:"modifiedBy" description:"修改人id"`
	ModifiedAt     int64                `json:"modifiedAt" description:"修改时间"`
}
