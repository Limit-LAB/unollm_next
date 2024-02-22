package httpHandler

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
)

const InjectedChatGLMHeader = "X-Inject-ChatGLM-Auth"
const InjectedChatGPTHeader = "X-Inject-ChatGPT-Auth"

type OpenAIEmbeddingRequest struct {
	Input          string `json:"input"`
	Model          string `json:"model"`
	EncodingFormap string `json:"encoding_format"`
}

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
		// TODO: Model Compatitable
		if strings.HasPrefix(req.Model, "chatglm") {
			ChatGLM_ChatCompletionHandler(c, opt.KeyTransformer, req)
			return
		}
		if strings.HasPrefix(req.Model, "chatgpt") {
			ChatGPT_ChatCompletitionsHandler(c, opt.KeyTransformer, req)
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

	r.POST("/embeddings", func(c *gin.Context) {
		var req OpenAIEmbeddingRequest
		err := c.BindJSON(&req)
		if err != nil {
			internalServerError(c, err)
			return
		}
		if strings.HasPrefix(req.Model, "chatglm::") {
			ChatGLM_EmbeddingHandler(c, req)
		}
		if strings.HasPrefix(req.Model, "chatgpt::") {
			ChatGPT_EmbeddingHandler(c, req)
		}
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
