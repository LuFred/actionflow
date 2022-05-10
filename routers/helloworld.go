package routers

import (
	"actionflow/core"
	"actionflow/dto"
	"actionflow/handlers"
)

func RegisterHelloworldHandlers() core.Route {
	return core.Route{
		Pattern: "/helloworld",
		Handler: &handlers.HelloworldHandler{},
		Routes: []core.RouteDescribe{
			{
				Verb:   core.Get,
				Path:   "/hi",
				Method: "Hi",
				Dto:    &dto.HelloworldRequest{},
			},
			{
				Verb:   core.Post,
				Path:   "/hii",
				Method: "Hi",
				Dto:    &dto.HelloworldRequest{},
			},
		},
	}

}
