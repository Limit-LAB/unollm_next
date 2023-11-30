package httpHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/provider/ChatGLM"
	"strings"
)

func NewOpenAIClient(c *gin.Context, tx KeyTransformer) *openai.Client {
	authHeader := getAuthHeader(c, InjectedChatGPTHeader)
	if authHeader == "" {
		return nil
	}

	endPoint := ""
	if tx != nil {
		rst, err := tx(authHeader, "chatgpt")
		if err != nil {
			return nil
		}
		if rst.EndPoint != "" {
			endPoint = rst.EndPoint
		}
		authHeader = rst.Key
	}
	if headerEp := c.GetHeader("X-OpenAI-Endpoint"); headerEp != "" {
		endPoint = headerEp
	}

	config := openai.DefaultConfig(authHeader)
	if endPoint != "" {
		config.BaseURL = endPoint
	}
	return openai.NewClientWithConfig(config)
}

func NewChatGLMClient(c *gin.Context, tx KeyTransformer) *ChatGLM.Client {
	authHeader := getAuthHeader(c, InjectedChatGLMHeader)
	if authHeader == "" {
		return nil
	}
	if tx == nil {
		return ChatGLM.NewClient(authHeader)
	}
	rst, err := tx(authHeader, "chatglm")
	if err != nil {
		return nil
	}
	cli := ChatGLM.NewClient(rst.Key)
	if rst.EndPoint != "" {
		cli.SetBase(rst.EndPoint)
	}
	return cli
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
