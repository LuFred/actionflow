package flowdto

type CreateJobRunInstanceRequest struct {
	JobId string `json:"jobId" transport:"path" required:"true" description:"flow job id"`
}
