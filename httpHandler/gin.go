package httpHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/relay"
)

const InjectedChatGLMHeader = "X-Inject-ChatGLM-Auth"
const InjectedChatGPTHeader = "X-Inject-ChatGPT-Auth"

func getProvider(m string) string {
	return "openai"
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
		if autoErr(c, c.BindJSON(&req)) {
			return
		}
		// TODO: Model Compatitable
		switch getProvider(req.Model) {
		case "openai":
			ChatGPT_ChatCompletitionsHandler(c, opt.KeyTransformer, req)
		case "chatglm":
			ChatGLM_ChatCompletionHandler(c, opt.KeyTransformer, req)
		}
	})
	r.POST("/completions", func(c *gin.Context) {
		var req openai.CompletionRequest
		if autoErr(c, c.BindJSON(&req)) {
			return
		}
		ChatGPT_CompletitionsHandler(c, opt.KeyTransformer, req)

	})

	r.POST("/embeddings", func(c *gin.Context) {
		var req relay.CommonEmbdReq
		if autoErr(c, c.BindJSON(&req)) {
			return
		}
		switch getProvider(req.Model) {
		case "openai":
			var _req openai.EmbeddingRequest
			if autoErr(c, c.BindJSON(&_req)) {
				return
			}
			ChatGPT_EmbeddingHandler(c, opt.KeyTransformer, _req)

		case "chatglm":
			ChatGLM_EmbeddingHandler(c, opt.KeyTransformer, req)
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
