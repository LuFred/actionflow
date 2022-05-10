package flow

import (
	"actionflow/core"
	"actionflow/db/entity"
	"actionflow/dto/flowdto"
	"actionflow/logic"
	"actionflow/pkg/ormutil"
	"actionflow/pkg/timeutil"
	"context"
	"fmt"
	"github.com/pborman/uuid"
	"go.uber.org/zap"
	"sync"
)

func (l *FlowLogic) GetActionsByFlowId(req *flowdto.GetActionsByFlowIdRequest) (*flowdto.GetActionsByFlowIdResponse, error) {
	resp := &flowdto.GetActionsByFlowIdResponse{}

	o := goorm.NewOrm()

	var list []entity.ActionEntity

	err := o.Select(&list, goorm.Cond{
		"flowId": req.FlowId,
	}, "createdAt")

	if err != nil {
		l.Logger.Error(err.Error())
		return nil, core.HttpInternalServerError()
	}

	resp.Data = make([]*flowdto.ActionDto, 0, len(list))

	if err != nil && err != goorm.ErrNoMoreRows {
		l.Logger.Error(err.Error())
		return nil, core.HttpInternalServerError()
	}

	for _, action := range list {
		resp.Data = append(resp.Data, convertActionEntityToDto(&action))
	}

	var wg sync.WaitGroup
	for _, dt := range resp.Data {
		wg.Add(1)
		go func(d *flowdto.ActionDto) {
			defer wg.Done()
			var edgeList []entity.FlowActionEdgeEntity
			goErr := o.Select(&edgeList, goorm.Cond{
				"endActionId": d.Id,
				"hops":        0,
			})

			if goErr != nil {
				l.Logger.Error(goErr.Error())
				err = core.HttpInternalServerError()
				return
			}

			d.PreIds = make([]string, 0, len(edgeList))
			for _, v := range edgeList {
				d.PreIds = append(d.PreIds, v.StartActionId)
			}
		}(dt)
	}
	wg.Wait()
	return resp, nil
}

func (l *FlowLogic) CreateAction(req *flowdto.CreateActionRequest) (*flowdto.ActionDto, error) {
	lg := l.Logger
	resp := &flowdto.ActionDto{}

	err := validateCreateActionRequest(req)
	if err != nil {
		return nil, core.HttpBadRequestCustomError(int32(logic.HttpErrorParameter), err.Error())
	}

	// check preIds
	invalidIds, err := checkActionsIsExists(l.ctx, req.PreIds, req.FlowId)
	if err != nil {
		if len(invalidIds) > 0 {
			return nil, core.HttpBadRequestError(logic.HttpErrorInvalidPreId)
		}

		lg.Error(err.Error())
		return nil, core.HttpInternalServerError()
	}

	//check nextIds
	invalidIds, err = checkActionsIsExists(l.ctx, req.NextIds, req.FlowId)
	if err != nil {
		if len(invalidIds) > 0 {
			return nil, core.HttpBadRequestError(logic.HttpErrorInvalidNextId)
		}

		lg.Error(err.Error())
		return nil, core.HttpInternalServerError()
	}

	actionEty := &entity.ActionEntity{}
	o := goorm.NewOrm()

	err = o.Tx(l.ctx, func(tx goorm.ITx) error {
		flowEty := &entity.FlowEntity{}
		// check flow is exists
		err := tx.One(flowEty, goorm.Cond{
			"id": req.FlowId,
		})

		if err != nil {
			if err == goorm.ErrNoMoreRows {
				return core.HttpBadRequestError(logic.HttpErrorFlowNotFound)
			}

			l.Logger.Error(err.Error(),
				zap.String("id", req.FlowId))
			return core.HttpInternalServerError()
		}
		err = tx.One(actionEty, goorm.Cond{
			"name":   req.Name,
			"flowId": req.FlowId,
		})

		if err == nil {
			return core.HttpBadRequestError(logic.HttpErrorNameAlreadyExists)
		}

		if err != goorm.ErrNoMoreRows {
			l.Logger.Error(err.Error(),
				zap.String("name", req.Name))
			return core.HttpInternalServerError()
		}

		curTime := timeutil.GetTimeMillisecond()
		actionEty.Id = uuid.New()
		actionEty.Name = req.Name
		actionEty.DisplayName = req.DisplayName
		actionEty.FlowId = req.FlowId
		actionEty.Type = req.Type
		actionEty.Command = req.Command
		actionEty.CreatedAt = curTime
		actionEty.CreatedBy = l.oauthInfo.UserId
		_, err = tx.Insert(actionEty)
		if err != nil {
			lg.Error(err.Error())
			return core.HttpInternalServerError()
		}

		//create action edge
		for _, startVertex := range req.PreIds {
			err = createActionEdage(tx, lg, startVertex, actionEty.Id)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		lg.Error(err.Error())
		return nil, err
	}

	resp = convertActionEntityToDto(actionEty)
	return resp, nil
}

func createActionEdage(tx goorm.ITx, lg *zap.Logger, startId string, endId string) error {
	edgeEty := &entity.FlowActionEdgeEntity{}
	//0.判断是否重复
	err := tx.One(edgeEty, goorm.Cond{
		"startActionId": startId,
		"endActionId":   endId,
		"hops":          0,
	})
	if err == nil {
		return core.HttpBadRequestError(logic.HttpErrorRepeatEdge)
	}

	if err != nil && err != goorm.ErrNoMoreRows {
		lg.Error(err.Error(),
			zap.String("startActionId", startId),
			zap.String("endActionId", endId))
		return core.HttpInternalServerError()
	}
	//1.判断是否存在环
	err = tx.One(edgeEty, goorm.Cond{
		"startActionId": endId,
		"endActionId":   startId,
	})
	if err == nil {
		return core.HttpBadRequestError(logic.HttpErrorGraphHaveCycle)
	}

	if err != nil && err != goorm.ErrNoMoreRows {
		lg.Error(err.Error(),
			zap.String("startActionId", startId),
			zap.String("endActionId", endId))
		return core.HttpInternalServerError()
	}

	//2.插入新边
	edgeEty.StartActionId = startId
	edgeEty.EndActionId = endId
	edgeEty.Hops = 0
	_id, err := tx.Insert(edgeEty)
	if err != nil {
		lg.Error(err.Error())
		return core.HttpInternalServerError()
	}
	id := _id.(int64)

	//更新新边
	edgeEty.Id = id
	edgeEty.EntryEdgeId = id
	edgeEty.ExitEdgeId = id
	edgeEty.DirectEdgeId = id

	err = tx.Update(edgeEty)
	if err != nil {
		lg.Error(err.Error())
		return core.HttpInternalServerError()
	}

	//3.新增start端到end端的边
	_, err = tx.Exec(`INSERT INTO flow_action_edge ( entryEdgeId, directEdgeId, exitEdgeId, startActionId, endActionId, hops ) SELECT
							id,
							?,
							?,
							startActionId,
							?,
							hops + 1 
							FROM
								flow_action_edge 
							WHERE
								endActionId = ?;`,
		id,
		id,
		endId,
		startId)

	if err != nil {
		lg.Error(err.Error())
		return core.HttpInternalServerError()
	}
	//4.新增start端到end端的出边
	_, err = tx.Exec(`INSERT INTO flow_action_edge ( entryEdgeId, directEdgeId, exitEdgeId, startActionId, endActionId, hops ) SELECT
							?,
							?,
							id,
							?,
							endActionId,
							hops + 1 
							FROM
								flow_action_edge 
							WHERE
								startActionId = ?;`,
		id,
		id,
		startId,
		endId)

	if err != nil {
		lg.Error(err.Error())
		return core.HttpInternalServerError()
	}
	//5.start的输入边到end的输出边的结束点
	_, err = tx.Exec(`INSERT INTO flow_action_edge ( entryEdgeId, directEdgeId, exitEdgeId, startActionId, endActionId, hops ) SELECT
							A.id,
							?,
							B.id,
							A.startActionId,
							B.endActionId,
							A.hops + B.hops + 1 
							FROM
								flow_action_edge A
								CROSS JOIN flow_action_edge B 
							WHERE
								A.endActionId = ? 
								AND B.startActionId = ?;`,
		id,
		startId,
		endId)

	if err != nil {
		lg.Error(err.Error())
		return core.HttpInternalServerError()
	}
	return nil
}

func validateCreateActionRequest(dto *flowdto.CreateActionRequest) error {
	if _, ok := ActionTypeId[dto.Type]; !ok {
		return fmt.Errorf("incorrect type; support:[blank, HTTP, Page]")
	}

	if len(dto.PreIds) == 0 {
		return fmt.Errorf(" preIds can not be null")
	}

	return nil
}

func checkActionsIsExists(ctx context.Context, ids []string, flowId string) (invalidIds []string, err error) {
	var wg sync.WaitGroup
	var mutex sync.Mutex
	o := goorm.NewOrm()
	for i := range ids {
		wg.Add(1)
		go func(index int32) {
			defer wg.Done()
			ety := &entity.ActionEntity{}
			er := o.One(ety, goorm.Cond{
				"id":     ids[index],
				"flowId": flowId,
			})

			if er != nil {
				err = er
				if er == goorm.ErrNoMoreRows {
					mutex.Lock()
					defer mutex.Unlock()
					invalidIds = append(invalidIds, ids[index])
				}
			}
		}(int32(i))
	}

	wg.Wait()
	return
}

func convertActionEntityToDto(in *entity.ActionEntity) *flowdto.ActionDto {
	rs := &flowdto.ActionDto{
		Id:          in.Id,
		FlowId:      in.FlowId,
		DisplayName: in.DisplayName,
		Name:        in.Name,
		Type:        in.Type,
		Command:     in.Command,
		CreatedBy:   in.CreatedBy,
		CreatedAt:   in.CreatedAt,
		ModifiedBy:  in.ModifiedBy,
		ModifiedAt:  in.ModifiedAt,
	}

	return rs
}
