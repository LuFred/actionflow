package appsflow

type (
	JobRunInstance struct {
		Id           string `json:"id" description:""`
		FlowId       string `json:"flowId" description:"流程id"`
		JobId        string `json:"jobId" description:"job id"`
		FlowVariable string `json:"flowVariable" description:"流程变量"`
		Timeout      int64  `json:"timeout" description:"超时时间"`
		Status       string `json:"status" description:"状态 [pending | executing | success | failed | cancelled]"`
		Message      string `json:"message" description:"结果消息"`
		ExecutedAt   int64  `json:"executedAt" description:"任务开始时间"`
		FinishedAt   int64  `json:"finishedAt" description:"任务结束时间"`
		CreatedBy    string `json:"createdBy" description:"创建人id"`
		CreatedAt    int64  `json:"createdAt" description:"创建时间"`
		ModifiedBy   string `json:"modifiedBy" description:"修改人id"`
		ModifiedAt   int64  `json:"modifiedAt" description:"修改时间"`
	}

	Job struct {
		Id           string `json:"id" description:""`
		FlowId       string `json:"flowId" description:"流程id"`
		FlowVariable string `json:"flowVariable" description:"流程变量"`
		Timeout      int64  `json:"timeout" description:"超时时间"`
		Status       string `json:"status" description:"任务状态"`
		Message      string `json:"message" description:"结果消息"`
		CreatedBy    string `json:"createdBy" description:"创建人id"`
		CreatedAt    int64  `json:"createdAt" description:"创建时间"`
		ModifiedBy   string `json:"modifiedBy" description:"修改人id"`
		ModifiedAt   int64  `json:"modifiedAt" description:"修改时间"`
	}

	Action struct {
		PreIds      []string `json:"preIds" description:"父节点id"`
		Id          string   `json:"id" description:""`
		FlowId      string   `json:"flowId" description:"流程id"`
		Name        string   `json:"name" description:"name"`
		DisplayName string   `json:"displayName" description:"显示名称"`
		Type        string   `json:"type" description:"节点类型[conditional | action]"`
		Command     string   `json:"command" description:"操作命令"`
		CreatedBy   string   `json:"createdBy" description:"创建人id"`
		CreatedAt   int64    `json:"createdAt" description:"创建时间"`
		ModifiedBy  string   `json:"modifiedBy" description:"修改人id"`
		ModifiedAt  int64    `json:"modifiedAt" description:"修改时间"`
	}

	FlowAction struct {
		Data []*Action `json:"data"`
	}
)
