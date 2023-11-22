package httpHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"strings"
)

func RegisterRoute(r *gin.Engine) {
	r.GET("/chat/completions", func(c *gin.Context) {
		var req openai.ChatCompletionRequest
		err := c.BindJSON(&req)
		if err != nil {
			// TODO: log
			return
		}
		if strings.HasPrefix(req.Model, "chatglm") {
			ChatGLM_ChatCompletionHandler(c, req)
			return
		}
		if strings.HasPrefix(req.Model, "chatgpt") {
			ChatGPT_ChatCompletitionsHandler(c, req)
			return
		}
	})
	r.GET("/completions", func(c *gin.Context) {
		var req openai.CompletionRequest
		err := c.BindJSON(&req)
		if err != nil {
		}
		if strings.HasPrefix(req.Model, "chatglm") {
			// TODO
		}
		if strings.HasPrefix(req.Model, "chatgpt") {

		}

	})
}
