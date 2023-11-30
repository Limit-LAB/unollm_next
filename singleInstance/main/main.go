package main

import (
	"github.com/KevinZonda/GoX/pkg/panicx"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.limit.dev/unollm/httpHandler"
	"go.limit.dev/unollm/singleInstance/shared"
	"go.limit.dev/unollm/singleInstance/svc"
	"go.limit.dev/unollm/singleInstance/svc/keyStore"
)

func main() {
	initAll()
	g := gin.New()
	g.Use(gin.Logger())
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AllowHeaders = append(corsCfg.AllowHeaders, "Authorization", "Authorisation")
	g.Use(cors.New(corsCfg))
	g.Use(gin.Recovery())

	httpHandler.RegisterRoute(g, httpHandler.RegisterOpt{
		KeyTransformer: keyStore.KeyTransformer,
	})

	registerRoute(g)

	err := g.Run(shared.GetCfg().ListenAddr)
	if err != nil {
		panic(err)
	}
}

func registerRoute(g gin.IRouter) {
	svc.RegisterSvc(g,
		map[string]svc.Svc{
			"keyStore/": &keyStore.KeyStoreSvc{},
		},
	)
}

func initAll() {
	err := shared.InitCfgModel("config.json")
	panicx.PanicIfNotNil(err, err)
	err = shared.InitMySql()
	panicx.PanicIfNotNil(err, err)
}
