package flowdto

type GetJobRunInstanceRequest struct {
	JobId         string `json:"jobId" transport:"path" required:"true" description:"job id"`
	RunInstanceId string `json:"runInstanceId" transport:"path" required:"true" description:"job runinstance id"`
}
