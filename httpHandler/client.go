package httpHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/provider/ChatGLM"
	"strings"
)

func NewOpenAIClient(c *gin.Context) *openai.Client {
	authHeader := getAuthHeader(c, InjectedChatGPTHeader)
	if authHeader == "" {
		return nil
	}

	config := openai.DefaultConfig(authHeader)
	endPoint := c.GetHeader("X-OpenAI-Endpoint")
	if endPoint != "" {
		config.BaseURL = endPoint
	}
	return openai.NewClientWithConfig(config)
}

func NewChatGLMClient(c *gin.Context) *ChatGLM.Client {
	authHeader := getAuthHeader(c, InjectedChatGLMHeader)
	if authHeader == "" {
		return nil
	}
	return ChatGLM.NewClient(authHeader)
}

func getAuthHeader(c *gin.Context, headers ...string) string {
	headers = append([]string{"Authorization", "Authorisation"}, headers...)
	authHeader := ""

	for _, header := range headers {
		authHeader = c.GetHeader(header)
		if authHeader != "" {
			break
		}
	}

	if strings.HasPrefix(authHeader, "Bearer ") {
		authHeader = authHeader[7:]
	}
	return authHeader
}
