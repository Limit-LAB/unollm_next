package user

import "github.com/gin-gonic/gin"

type UserSvc struct{}

func (svc *UserSvc) RegisterRouter(g gin.IRouter) {
	g.POST("/login", login)
	g.POST("/create", create)
	g.POST("/logout", logout)
}
