package handlers

import (
	"actionflow/dto/flowdto"
	"actionflow/logic/flow"
	"actionflow/logic/task"
)

// CreateFlow
func (h *ActionFlowHandler) CreateFlow(body *flowdto.CreateFlowRequest) {
	l := flow.NewFlowLogic(h.Ctx)
	resp, err := l.CreateFlow(body)
	h.Next(resp, err)
}

// CreateAction
func (h *ActionFlowHandler) CreateAction(body *flowdto.CreateActionRequest) {
	l := flow.NewFlowLogic(h.Ctx)
	resp, err := l.CreateAction(body)
	h.Next(resp, err)
}

// GetActionsByFlowId
func (h *ActionFlowHandler) GetActionsByFlowId(body *flowdto.GetActionsByFlowIdRequest) {
	l := flow.NewFlowLogic(h.Ctx)
	resp, err := l.GetActionsByFlowId(body)
	h.Next(resp, err)
}

// CreateJob
func (h *ActionFlowHandler) CreateJob(body *flowdto.CreateFlowJobRequest) {
	l := task.NewTaskLogic(h.Ctx)
	resp, err := l.CreateJob(body)
	h.Next(resp, err)
}

// GetJobById
func (h *ActionFlowHandler) GetJobById(body *flowdto.GetJobByIdRequest) {
	l := task.NewTaskLogic(h.Ctx)
	resp, err := l.GetJobById(body)
	h.Next(resp, err)
}

// CreateJobRunInstance
func (h *ActionFlowHandler) CreateJobRunInstance(body *flowdto.CreateJobRunInstanceRequest) {
	l := task.NewTaskLogic(h.Ctx)
	resp, err := l.CreateJobRunInstance(body)
	h.Next(resp, err)
}

// GetJobRunInstanceById
func (h *ActionFlowHandler) GetJobRunInstanceById(body *flowdto.GetJobRunInstanceRequest) {
	l := task.NewTaskLogic(h.Ctx)
	resp, err := l.GetJobRunInstanceById(body)
	h.Next(resp, err)
}
