package httpHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"limit.dev/unollm/provider/ChatGLM"
	"limit.dev/unollm/relay"
	"limit.dev/unollm/utils"
	"log"
	"net/http"
)

func ChatGLM_ChatCompletionHandler(c *gin.Context, req openai.ChatCompletionRequest) {
	auth := utils.GetAuthorisation(c)
	if auth == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "no api key provided",
		})
		c.Abort()
		return
	}
	zpReq := ChatGLM.ChatCompletionRequest{
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Incremental: req.Stream,
	}

	for _, m := range req.Messages {
		zpReq.Prompt = append(zpReq.Prompt, ChatGLM.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}
	cli := ChatGLM.NewClient(auth)
	if !req.Stream {
		rst, err := cli.ChatCompletion(zpReq, req.Model)
		if err != nil {
			log.Println(err)
			return
		}
		c.JSON(200, rst)
	}
	llm, result, err := cli.ChatCompletionStreamingRequest(zpReq, req.Model)
	defer safeClose(llm)
	defer safeClose(result)

	if err != nil {
		log.Println(err)
		return
	}
	relay.ChatGlmStream2OpenAI(c, llm, result)
}

func safeClose[T any](ch chan<- T) {
	defer func() { recover() }()
	close(ch)
}
