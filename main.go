package main

import (
	"go.limit.dev/unollm/relay/respTransformer"
	"log"
	"net"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/ChatGLM"
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
		g.POST("/v1/chat/completions", func(c *gin.Context) {
			var req openai.ChatCompletionRequest
			err := c.BindJSON(&req)
			if err != nil {
				log.Println(err)
				return
			}
			cli := ChatGLM.NewClient(zhipuaiApiKey)

			zpReq := ChatGLM.ChatCompletionRequest{
				Temperature: req.Temperature,
				TopP:        req.TopP,
				Incremental: true,
			}

			for _, m := range req.Messages {
				zpReq.Prompt = append(zpReq.Prompt, ChatGLM.ChatCompletionMessage{
					Role:    m.Role,
					Content: m.Content,
				})
			}
			resp, err := cli.ChatCompletionStreamingRequest(zpReq, "chatglm_turbo")
			if err != nil {
				log.Println(err)
				return
			}
			respTransformer.ChatGLMToOpenAIStream(c, resp)
		})
		g.Run("127.0.0.1:11451")
	}()

	// start grpc server
	grpcServer := grpc.NewServer()
	lis, err := net.Listen("tcp", "127.0.0.1:19198")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	service := grpcServer.UnoForwardServer{}
	model.RegisterUnoLLMv1Server(grpcServer, &service)
	grpcServer.Serve(lis)
}
