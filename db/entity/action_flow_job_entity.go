package entity

type ActionFlowJobEntity struct {
	Id             string `db:"id" description:""`
	FlowId         string `db:"flowId" description:"actionflow id"`
	FlowVariable   string `db:"flowVariable" description:"action flow variable"`
	Timeout        int64  `db:"timeout" description:"timeout"`
	Status         string `db:"status" description:"job status [deleted | pause | active]"`
	Message        string `db:"message" description:"result message"`
	RetryOptionsId string `db:"retryOptionsId" description:"retry options id"`
	CreatedBy      string `db:"createdBy" description:"created user id"`
	CreatedAt      int64  `db:"createdAt" description:"created time"`
	ModifiedBy     string `db:"modifiedBy" description:"modify user id"`
	ModifiedAt     int64  `db:"modifiedAt" description:"modified time"`
}

func (d *ActionFlowJobEntity) TableName() string {
	return "action_flow_job"
}
