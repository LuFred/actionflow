package dto

type RetryOptionsDto struct {
	Id                       []byte `json:"id" description:""`
	ResourceType             string `json:"resourceType" description:"目标资源类型"`
	ResourceId               string `json:"resourceId" description:"目标资源id"`
	Attempts                 int    `json:"attempts" description:"当前尝试次数"`
	InitialInterval          int    `json:"initialInterval" description:"重试默认时间间隔"`
	Coefficient              int    `json:"coefficient" description:"重试时间系数"`
	MaximumInterval          int    `json:"maximumInterval" description:"最大时间间隔"`
	MaximumAttempts          int    `json:"maximumAttempts" description:"最大尝试次数"`
	NonRetryableErrorReasons int    `json:"nonRetryableErrorReasons" description:"不应重试的错误"`
	CreatedBy                []byte `json:"createdBy" description:"创建人id"`
	CreatedAt                int64  `json:"createdAt" description:"创建时间"`
	ModifiedBy               []byte `json:"modifiedBy" description:"修改人id"`
	ModifiedAt               int64  `json:"modifiedAt" description:"修改时间"`
}
