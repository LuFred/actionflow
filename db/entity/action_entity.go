package entity

type ActionEntity struct {
	Id             string `db:"id" description:""`
	DisplayName    string `db:"displayName" description:"显示名称"`
	Name           string `db:"name" description:"调用名称"`
	FlowId         string `db:"flowId" description:"关联小程序id"`
	Type           string `db:"type" description:"动作类型"`
	Command        string `db:"command" description:"动作描述"`
	RetryOptionsId string `db:"retryOptionsId" description:"重试策略"`
	CreatedBy      string `db:"createdBy" description:"创建人id"`
	CreatedAt      int64  `db:"createdAt" description:"创建时间"`
	ModifiedBy     string `db:"modifiedBy" description:"修改人id"`
	ModifiedAt     int64  `db:"modifiedAt" description:"修改时间"`
}

func (d *ActionEntity) TableName() string {
	return "action"
}
