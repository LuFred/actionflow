package flowdto

import "actionflow/dto"

type FlowDto struct {
	Id              string               `json:"id" description:""`
	AppId           string               `json:"appId" description:"关联小程序id"`
	DisplayName     string               `json:"displayName" description:"显示名称"`
	Name            string               `json:"name" description:"调用名称"`
	Description     string               `json:"description" description:"描述"`
	InputParameter  string               `json:"inputParameter" description:"输入参数"`
	OutputParameter string               `json:"outputParameter" description:"输出参数"`
	Timeout         int64                `json:"timeout" description:"超时时间(s)"`
	RetryOptionsId  *dto.RetryOptionsDto `json:"retryOptionsId" description:"重试策略"`
	CreatedBy       string               `json:"createdBy" description:"创建人id"`
	CreatedAt       int64                `json:"createdAt" description:"创建时间"`
	ModifiedBy      string               `json:"modifiedBy" description:"修改人id"`
	ModifiedAt      int64                `json:"modifiedAt" description:"修改时间"`
}
