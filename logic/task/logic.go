package task

import (
	"actionflow/core"
	"context"
	"go.uber.org/zap"
)

type TaskLogic struct {
	Logger    *zap.Logger
	ctx       context.Context
	srvCtx    *core.Server
	oauthInfo *core.OauthInfo
}

func NewTaskLogic(ctx *core.Context) TaskLogic {
	return TaskLogic{
		Logger:    ctx.SrvCtx.Logger(),
		ctx:       ctx.Request.Context(),
		srvCtx:    ctx.SrvCtx,
		oauthInfo: ctx.OauthInfo,
	}
}
