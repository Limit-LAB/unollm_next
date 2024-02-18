package httpHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"strings"
)

const InjectedChatGLMHeader = "X-Inject-ChatGLM-Auth"
const InjectedChatGPTHeader = "X-Inject-ChatGPT-Auth"

func RegisterRoute(r *gin.Engine, opt RegisterOpt) {
	if opt.ChatGPTKey != "" {
		r.Use(func(c *gin.Context) {
			c.Request.Header.Set(InjectedChatGPTHeader, opt.ChatGPTKey)
		})
	}
	if opt.ChatGLMKey != "" {
		r.Use(func(c *gin.Context) {
			c.Request.Header.Set(InjectedChatGLMHeader, opt.ChatGLMKey)
		})
	}

	r.POST("/chat/completions", func(c *gin.Context) {
		var req openai.ChatCompletionRequest
		err := c.BindJSON(&req)
		if err != nil {
			internalServerError(c, err)
			return
		}
		if strings.HasPrefix(req.Model, "chatglm") {
			ChatGLM_ChatCompletionHandler(c, opt.KeyTransformer, req)
			return
		}
		if strings.HasPrefix(req.Model, "chatgpt") {
			ChatGPT_ChatCompletitionsHandler(c, opt.KeyTransformer, req)
			return
		}
	})
	r.POST("/completions", func(c *gin.Context) {
		var req openai.CompletionRequest
		err := c.BindJSON(&req)
		if err != nil {
			internalServerError(c, err)
			return
		}
		ChatGPT_CompletitionsHandler(c, opt.KeyTransformer, req)

	})
}

type RegisterOpt struct {
	ChatGLMKey     string
	ChatGPTKey     string
	KeyTransformer KeyTransformer
}

type KeyTransformerResult struct {
	Key      string
	EndPoint string
}

type KeyTransformer func(key string, provider string) (KeyTransformerResult, error)
