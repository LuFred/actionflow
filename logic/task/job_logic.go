package task

import (
	"actionflow/core"
	"actionflow/db/entity"
	"actionflow/dto/flowdto"
	"actionflow/logic"
	goorm "actionflow/pkg/ormutil"
	"actionflow/pkg/timeutil"
	"github.com/pborman/uuid"
	"go.uber.org/zap"
)

func (l *TaskLogic) CreateJob(req *flowdto.CreateFlowJobRequest) (*flowdto.FlowJobDto, error) {
	lg := l.Logger
	resp := &flowdto.FlowJobDto{}
	jobEty := &entity.ActionFlowJobEntity{}
	o := goorm.NewOrm()
	err := o.Tx(l.ctx, func(tx goorm.ITx) error {
		flowEty := &entity.FlowEntity{}
		err := tx.One(flowEty, goorm.Cond{
			"id": req.FlowId,
		})

		if err != nil {
			if err == goorm.ErrNoMoreRows {
				return core.HttpBadRequestError(logic.HttpErrorFlowNotFound)
			}

			l.Logger.Error(err.Error(),
				zap.String("flowId", req.FlowId))
			return core.HttpInternalServerError()
		}

		jobEty.Id = uuid.New()
		jobEty.FlowId = req.FlowId
		jobEty.FlowVariable = req.FlowVariable
		jobEty.Timeout = req.Timeout
		jobEty.Status = JobStatusTypePause.String()
		jobEty.CreatedAt = timeutil.GetTimeMillisecond()
		jobEty.CreatedBy = l.oauthInfo.UserId
		_, err = tx.Insert(jobEty)

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

	resp = convertFlowJobEntityToDto(jobEty)
	return resp, nil
}

func (l *TaskLogic) GetJobById(req *flowdto.GetJobByIdRequest) (*flowdto.FlowJobDto, error) {
	resp := &flowdto.FlowJobDto{}
	jobEty := &entity.ActionFlowJobEntity{}
	o := goorm.NewOrm()
	err := o.One(jobEty, goorm.Cond{
		"id": req.JobId,
	})
	if err != nil {
		if err == goorm.ErrNoMoreRows {
			return nil, core.HttpNotFoundError()
		}
		l.Logger.Error(err.Error(),
			zap.String("jobId", req.JobId))
		return nil, core.HttpInternalServerError()
	}

	resp = convertFlowJobEntityToDto(jobEty)
	return resp, nil
}

func convertFlowJobEntityToDto(in *entity.ActionFlowJobEntity) *flowdto.FlowJobDto {
	rs := &flowdto.FlowJobDto{
		Id:           in.Id,
		FlowId:       in.FlowId,
		FlowVariable: in.FlowVariable,
		Timeout:      in.Timeout,
		Status:       in.Status,
		Message:      in.Message,
		CreatedBy:    in.CreatedBy,
		CreatedAt:    in.CreatedAt,
		ModifiedBy:   in.ModifiedBy,
		ModifiedAt:   in.ModifiedAt,
	}
	return rs
}
