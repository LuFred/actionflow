package entity

type FlowEntity struct {
	Id              string `db:"id" description:""`
	AppId           string `db:"appId" description:"关联小程序id"`
	DisplayName     string `db:"displayName" description:"显示名称"`
	Name            string `db:"name" description:"调用名称"`
	Description     string `db:"description" description:"描述"`
	InputParameter  string `db:"inputParameter" description:"输入参数"`
	OutputParameter string `db:"outputParameter" description:"输出参数"`
	Timeout         int64  `db:"timeout" description:"超时时间(s)"`
	RetryOptionsId  string `db:"retryOptionsId" description:"重试策略"`
	CreatedBy       string `db:"createdBy" description:"创建人id"`
	CreatedAt       int64  `db:"createdAt" description:"创建时间"`
	ModifiedBy      string `db:"modifiedBy" description:"修改人id"`
	ModifiedAt      int64  `db:"modifiedAt" description:"修改时间"`
}

func (d *FlowEntity) TableName() string {
	return "action_flow"
}
