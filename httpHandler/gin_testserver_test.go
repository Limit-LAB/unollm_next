package httpHandler_test

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/httpHandler"
)

var GetClient = testServer()

func testServer() (client func(key string) *openai.Client) {
	g := gin.New()
	g.Use(gin.Logger())
	corsCfg := cors.DefaultConfig()
	corsCfg.AllowAllOrigins = true
	corsCfg.AllowHeaders = append(corsCfg.AllowHeaders, "Authorization", "Authorisation")
	g.Use(cors.New(corsCfg))
	g.Use(gin.Recovery())
	httpHandler.RegisterRoute(g, httpHandler.RegisterOpt{})
	addr := "127.0.0.1:11451"

	client = func(key string) (client *openai.Client) {
		config := openai.DefaultConfig(key)
		config.BaseURL = "http://" + addr
		client = openai.NewClientWithConfig(config)
		return
	}

	go func() {
		g.Run(addr)
	}()
	time.Sleep(3 * time.Second) // Wait for server to start
	return
}
