package main

import (
	"go.limit.dev/unollm/grpcServer"
	"go.limit.dev/unollm/httpHandler"
	"log"
	"net"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.limit.dev/unollm/model"
	"google.golang.org/grpc"
)

func main() {
	// start openai style server
	go func() {
		logger := logrus.New()
		logger.Out = os.Stdout
		logger.SetLevel(logrus.InfoLevel)

		godotenv.Load("./.env")
		g := gin.New()
		g.Use(gin.Logger())
		corsCfg := cors.DefaultConfig()
		corsCfg.AllowAllOrigins = true
		corsCfg.AllowHeaders = append(corsCfg.AllowHeaders, "Authorization")
		g.Use(cors.New(corsCfg))
		g.Use(gin.Recovery())

		zhipuaiApiKey := os.Getenv("TEST_ZHIPUAI_API")
		httpHandler.RegisterRoute(g, httpHandler.RegisterOpt{
			InjectChatGLMKey: true,
			ChatGLMKey:       zhipuaiApiKey,
		})
		g.Run("127.0.0.1:11451")
	}()

	// start grpc server
	svr := grpc.NewServer()
	tcpList, err := net.Listen("tcp", "127.0.0.1:19198")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	service := grpcServer.UnoForwardServer{}
	model.RegisterUnoLLMv1Server(svr, &service)
	svr.Serve(tcpList)
}
