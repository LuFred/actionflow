package flow

import (
	"actionflow/core"
	"actionflow/db/entity"
	"actionflow/dto/flowdto"
	"actionflow/logic"
	goorm "actionflow/pkg/ormutil"
	"actionflow/pkg/timeutil"
	"fmt"
	"github.com/pborman/uuid"
	"go.uber.org/zap"
)

func (l *FlowLogic) CreateFlow(req *flowdto.CreateFlowRequest) (*flowdto.FlowDto, error) {
	lg := l.Logger
	lg.Debug(fmt.Sprintf("oauthInfo %v", l.oauthInfo))
	resp := &flowdto.FlowDto{
		DisplayName: l.oauthInfo.UserId,
	}
	flowEty := &entity.FlowEntity{}
	// check same name
	o := goorm.NewOrm()
	err := o.Tx(l.ctx, func(tx goorm.ITx) error {
		err := tx.One(flowEty, goorm.Cond{
			"name":  req.Name,
			"appId": req.AppId,
		})

		if err == nil {
			return core.HttpBadRequestError(logic.HttpErrorNameAlreadyExists)
		}

		if err != goorm.ErrNoMoreRows {
			l.Logger.Error(err.Error(),
				zap.String("name", req.Name))
			return core.HttpInternalServerError()
		}

		//create action flow
		curTime := timeutil.GetTimeMillisecond()
		flowEty.Id = uuid.New()
		flowEty.Name = req.Name
		flowEty.DisplayName = req.DisplayName
		flowEty.AppId = req.AppId
		flowEty.Timeout = req.Timeout
		flowEty.Description = req.Description
		flowEty.InputParameter = req.InputParameter
		flowEty.OutputParameter = req.OutputParameter
		flowEty.CreatedAt = curTime
		flowEty.CreatedBy = l.oauthInfo.UserId

		_, err = tx.Insert(flowEty)
		if err != nil {
			lg.Error(err.Error())
			return core.HttpInternalServerError()
		}

		// create start action node
		actionEty := &entity.ActionEntity{}
		actionEty.Id = flowEty.Id
		actionEty.Name = "start"
		actionEty.DisplayName = "start"
		actionEty.FlowId = flowEty.Id
		actionEty.Type = ActionTypeBlank.String()
		actionEty.CreatedAt = curTime
		actionEty.CreatedBy = l.oauthInfo.UserId
		_, err = tx.Insert(actionEty)

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

	resp = convertActionFlowEntityToDto(flowEty)
	return resp, nil
}

func convertActionFlowEntityToDto(in *entity.FlowEntity) *flowdto.FlowDto {
	rs := &flowdto.FlowDto{
		Id:              in.Id,
		AppId:           in.AppId,
		DisplayName:     in.DisplayName,
		Name:            in.Name,
		Description:     in.Description,
		InputParameter:  in.InputParameter,
		OutputParameter: in.OutputParameter,
		Timeout:         in.Timeout,
		CreatedBy:       in.CreatedBy,
		CreatedAt:       in.CreatedAt,
		ModifiedBy:      in.ModifiedBy,
		ModifiedAt:      in.ModifiedAt,
	}

	return rs
}
