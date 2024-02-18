package relay_test

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/provider/ChatGLM"
	"go.limit.dev/unollm/relay/respTransformer"
)

func mockServer() *gin.Engine {
	godotenv.Load("../.env")
	g := gin.New()
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
			Model:       req.Model,
			Temperature: req.Temperature,
			TopP:        req.TopP,
		}

		for _, m := range req.Messages {
			zpReq.Messages = append(zpReq.Messages, ChatGLM.ChatCompletionMessage{
				Role:    m.Role,
				Content: m.Content,
			})
		}
		_r, err := cli.ChatCompletionStreamingRequest(zpReq)
		if err != nil {
			log.Println(err)
			return
		}
		respTransformer.ChatGLMToOpenAIStream(c, _r)
	})
	return g
}

func TestZhipuChatCompletionStream(t *testing.T) {
	addr := "127.0.0.1:11451"
	go func() {
		e := mockServer()
		e.Run(addr)
	}()
	time.Sleep(3 * time.Second) // Wait for server to start
	config := openai.DefaultConfig("114514")
	config.BaseURL = "http://" + addr + "/v1"
	client := openai.NewClientWithConfig(config)
	resp, err := client.CreateChatCompletionStream(context.Background(),
		openai.ChatCompletionRequest{
			Model: ChatGLM.ModelTurbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: "假如今天下大雨，我是否需要带伞？",
				},
			},
		})
	if err != nil {
		t.Error(err)
	}
	for {
		cv, e := resp.Recv()
		if e != nil {
			if errors.Is(e, io.EOF) {
				break
			}
			t.Error(e)
			break
		}
		log.Println(cv.Choices[0].Delta.Content)
	}
}
