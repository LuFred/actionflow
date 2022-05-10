package flow

import (
	"actionflow/core"
	"context"
	"go.uber.org/zap"
)

type FlowLogic struct {
	Logger    *zap.Logger
	ctx       context.Context
	srvCtx    *core.Server
	oauthInfo *core.OauthInfo
}

func NewFlowLogic(ctx *core.Context) FlowLogic {
	return FlowLogic{
		Logger:    ctx.SrvCtx.Logger(),
		ctx:       ctx.Request.Context(),
		srvCtx:    ctx.SrvCtx,
		oauthInfo: ctx.OauthInfo,
	}
}
