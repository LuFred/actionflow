package handlers

import (
	"actionflow/core"
	"actionflow/dto"
	"actionflow/logic"
)

//HelloworldHandler
type HelloworldHandler struct {
	core.Handler
}

// helloworld
func (h *HelloworldHandler) Hi(body *dto.HelloworldRequest) {
	l := logic.NewHelloworldLogic(h.Ctx)
	resp, err := l.Hi(body)
	h.Next(resp, err)
}
