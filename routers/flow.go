package routers

import (
	"actionflow/core"
	"actionflow/dto/flowdto"
	"actionflow/handlers"
)

func RegisterFlowHandlers() core.Route {
	return core.Route{
		Pattern: "/flows",
		Handler: &handlers.ActionFlowHandler{},
		Routes: []core.RouteDescribe{
			{
				Verb:   core.Post,
				Path:   "/",
				Method: "CreateFlow",
				Dto:    &flowdto.CreateFlowRequest{},
			},
			{
				Verb:   core.Post,
				Path:   "/actions",
				Method: "CreateAction",
				Dto:    &flowdto.CreateActionRequest{},
			},
			{
				Verb:   core.Get,
				Path:   "/{flowId}/actions",
				Method: "GetActionsByFlowId",
				Dto:    &flowdto.GetActionsByFlowIdRequest{},
			},
			{
				Verb:   core.Post,
				Path:   "/{flowId}/jobs",
				Method: "CreateJob",
				Dto:    &flowdto.CreateFlowJobRequest{},
			},
			{
				Verb:   core.Get,
				Path:   "/jobs/{jobId}",
				Method: "GetJobById",
				Dto:    &flowdto.GetJobByIdRequest{},
			},
			{
				Verb:   core.Post,
				Path:   "/jobs/{jobId}/runinstance",
				Method: "CreateJobRunInstance",
				Dto:    &flowdto.CreateJobRunInstanceRequest{},
			},
			{
				Verb:   core.Get,
				Path:   "/jobs/{jobId}/runinstance/{runInstanceId}",
				Method: "GetJobRunInstanceById",
				Dto:    &flowdto.GetJobRunInstanceRequest{},
			},
		},
	}
}
