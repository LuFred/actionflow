package task

import (
	"actionflow/core"
	"actionflow/db/entity"
	"actionflow/dto/flowdto"
	"actionflow/logic"
	goorm "actionflow/pkg/ormutil"
	"actionflow/pkg/temporalutil"
	"actionflow/pkg/timeutil"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

func (l *TaskLogic) GetJobRunInstanceById(req *flowdto.GetJobRunInstanceRequest) (*flowdto.FlowJobRunInstanceDto, error) {
	resp := &flowdto.FlowJobRunInstanceDto{}
	runInstanceEty := &entity.ActionFlowJobRunInstanceEntity{}
	o := goorm.NewOrm()
	err := o.One(runInstanceEty, goorm.Cond{
		"id": req.RunInstanceId,
	})

	if err != nil {
		if err == goorm.ErrNoMoreRows {
			return nil, core.HttpNotFoundError()
		}
		l.Logger.Error(err.Error(),
			zap.String("runInstanceId", req.RunInstanceId))
		return nil, core.HttpInternalServerError()
	}
	if runInstanceEty.JobId != req.JobId {
		return nil, core.HttpNotFoundError()
	}

	resp = convertFlowJobRunInstanceEntityToDto(runInstanceEty)
	return resp, nil
}

func (l *TaskLogic) CreateJobRunInstance(req *flowdto.CreateJobRunInstanceRequest) (*flowdto.FlowJobRunInstanceDto, error) {
	lg := l.Logger
	resp := &flowdto.FlowJobRunInstanceDto{}
	runInstanceEty := &entity.ActionFlowJobRunInstanceEntity{}
	o := goorm.NewOrm()
	err := o.Tx(l.ctx, func(tx goorm.ITx) error {
		jobEty := &entity.ActionFlowJobEntity{}
		err := tx.One(jobEty, goorm.Cond{
			"id": req.JobId,
		})

		if err != nil {
			if err == goorm.ErrNoMoreRows {
				return core.HttpBadRequestError(logic.HttpErrorFlowJobNotFound)
			}

			l.Logger.Error(err.Error(),
				zap.String("jobId", req.JobId))
			return core.HttpInternalServerError()
		}

		runInstanceEty.Id = uuid.New()
		runInstanceEty.FlowId = jobEty.FlowId
		runInstanceEty.JobId = req.JobId
		runInstanceEty.FlowVariable = jobEty.FlowVariable
		runInstanceEty.Timeout = jobEty.Timeout
		runInstanceEty.Status = JobRunInstanceStatusTypePending.String()
		runInstanceEty.CreatedAt = timeutil.GetTimeMillisecond()
		runInstanceEty.CreatedBy = l.oauthInfo.UserId
		_, err = tx.Insert(runInstanceEty)

		if err != nil {
			lg.Error(err.Error())
			return core.HttpInternalServerError()
		}

		//trigger run jobruninstance
		err = l.executeWorkflow(err, runInstanceEty.Id, runInstanceEty.JobId)
		if err != nil {
			lg.Error(err.Error())
			return core.HttpInternalServerError()
		}

		return nil
	})

	if err != nil {
		lg.Error(err.Error())
		return nil, err
	}

	resp = convertFlowJobRunInstanceEntityToDto(runInstanceEty)
	return resp, nil
}

func (l *TaskLogic) executeWorkflow(err error, runInstanceId string, jobId string) error {
	// call temporal to create jobRunInstance async
	c, err := temporalutil.GetClient()
	if err != nil {
		l.Logger.Error(err.Error())
		return core.HttpInternalServerError()
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        runInstanceId,
		TaskQueue: l.srvCtx.Cfg.TemporalConf.JobRunInstanceQueue,
	}
	we, err := c.ExecuteWorkflow(l.ctx, workflowOptions, l.srvCtx.Cfg.TemporalConf.JobRunInstanceWorkflow, runInstanceId, jobId)

	if err != nil {
		l.Logger.Error(err.Error())
		return core.HttpInternalServerError()
	}

	l.Logger.Info("started workflow", zap.String("workflowId", we.GetID()), zap.String("RunId", we.GetRunID()))
	return nil
}

func convertFlowJobRunInstanceEntityToDto(in *entity.ActionFlowJobRunInstanceEntity) *flowdto.FlowJobRunInstanceDto {
	rs := &flowdto.FlowJobRunInstanceDto{
		Id:           in.Id,
		FlowId:       in.FlowId,
		JobId:        in.JobId,
		FlowVariable: in.FlowVariable,
		Timeout:      in.Timeout,
		Status:       in.Status,
		Message:      in.Message,
		ExecutedAt:   in.ExecutedAt,
		FinishedAt:   in.FinishedAt,
		CreatedBy:    in.CreatedBy,
		CreatedAt:    in.CreatedAt,
		ModifiedBy:   in.ModifiedBy,
		ModifiedAt:   in.ModifiedAt,
	}

	return rs
}
