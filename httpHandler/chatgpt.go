package httpHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

func ChatGPT_ChatCompletitionsHandler(c *gin.Context, req openai.ChatCompletionRequest) {
	if req.Stream {
		return
	}
	config := openai.DefaultConfig(req.Model)
	endPoint := c.GetHeader("X-OpenAI-Endpoint")
	if endPoint != "" {
		config.BaseURL = endPoint
	}
	client := openai.NewClientWithConfig(config)
	resp, err := client.CreateChatCompletion(c, req)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, resp)
}
