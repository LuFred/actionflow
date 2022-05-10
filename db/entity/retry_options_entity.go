package entity

type RetryOptionsEntity struct {
	Id                       []byte `db:"id" description:""`
	ResourceType             string `db:"resourceType" description:"目标资源类型"`
	ResourceId               string `db:"resourceId" description:"目标资源id"`
	Attempts                 int    `db:"attempts" description:"当前尝试次数"`
	InitialInterval          int    `db:"initialInterval" description:"重试默认时间间隔"`
	Coefficient              int    `db:"coefficient" description:"重试时间系数"`
	MaximumInterval          int    `db:"maximumInterval" description:"最大时间间隔"`
	MaximumAttempts          int    `db:"maximumAttempts" description:"最大尝试次数"`
	NonRetryableErrorReasons int    `db:"nonRetryableErrorReasons" description:"不应重试的错误"`
	CreatedBy                []byte `db:"createdBy" description:"创建人id"`
	CreatedAt                int64  `db:"createdAt" description:"创建时间"`
	ModifiedBy               []byte `db:"modifiedBy" description:"修改人id"`
	ModifiedAt               int64  `db:"modifiedAt" description:"修改时间"`
}

func (d *RetryOptionsEntity) TableName() string {
	return "retry_options"
}
