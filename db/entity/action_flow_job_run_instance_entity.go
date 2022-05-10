package entity

type ActionFlowJobRunInstanceEntity struct {
	Id             string `db:"id" description:""`
	JobId          string `db:"jobId" description:"job id"`
	FlowId         string `db:"flowId" description:"job id"`
	FlowVariable   string `db:"flowVariable" description:"action flow variable"`
	Timeout        int64  `db:"timeout" description:"timeout"`
	Status         string `db:"status" description:"task status [pending | executing | success | failed | cancelled]"`
	Message        string `db:"message" description:"result message"`
	RetryOptionsId string `db:"retryOptionsId" description:"retry options id"`
	ExecutedAt     int64  `db:"executedAt" description:"executed time"`
	FinishedAt     int64  `db:"finishedAt" description:"finished time"`
	CreatedBy      string `db:"createdBy" description:"created user id"`
	CreatedAt      int64  `db:"createdAt" description:"created time"`
	ModifiedBy     string `db:"modifiedBy" description:"modify user id"`
	ModifiedAt     int64  `db:"modifiedAt" description:"modified time"`
}

func (d *ActionFlowJobRunInstanceEntity) TableName() string {
	return "action_flow_job_run_instance"
}
