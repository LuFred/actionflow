package flowdto

type CreateFlowRequest struct {
	Name            string `json:"name" transport:"body" required:"true" description:"name"`
	DisplayName     string `json:"displayName" transport:"body" required:"true" description:"显示名称"`
	AppId           string `json:"appId" transport:"body"  description:"app id"`
	Description     string `json:"description" transport:"body"  description:"描述"`
	InputParameter  string `json:"inputParameter" transport:"body" description:"输入参数"`
	OutputParameter string `json:"outputParameter" transport:"body"  description:"输出参数"`
	Timeout         int64  `json:"timeout" transport:"body"  description:"超时时间(s)"`
}
