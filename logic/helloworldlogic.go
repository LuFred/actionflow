package logic

import (
	"actionflow/core"
	"actionflow/dto"
	"context"
	"fmt"

	"go.uber.org/zap"
)

type HelloworldLogic struct {
	Logger *zap.Logger
	ctx    context.Context
	srvCtx *core.Server
}

func NewHelloworldLogic(ctx *core.Context) HelloworldLogic {
	return HelloworldLogic{
		Logger: ctx.SrvCtx.Logger(),
		ctx:    ctx.Request.Context(),
		srvCtx: ctx.SrvCtx,
	}
}

func (l *HelloworldLogic) Hi(req *dto.HelloworldRequest) (*dto.HelloworldResponse, error) {
	resp := &dto.HelloworldResponse{
		Message: fmt.Sprintf("%s %s", req.Name, l.srvCtx.Cfg.Port),
	}

	return resp, nil
}
