package flowdto

type GetJobByIdRequest struct {
	JobId string `json:"jobId" transport:"path" required:"true" description:"job id"`
}
