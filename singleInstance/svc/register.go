package svc

import "github.com/gin-gonic/gin"

type Svc interface {
	RegisterRouter(g gin.IRouter)
}

func RegisterSvc(g gin.IRouter, svcs map[string]Svc) {
	for key, svc := range svcs {
		if svc == nil {
			continue
		}
		ng := g.Group(key)
		svc.RegisterRouter(ng)
	}
}
