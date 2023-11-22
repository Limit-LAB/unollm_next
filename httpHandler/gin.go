package httpHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"limit.dev/unollm/provider/ChatGLM"
)


func router(r *gin.Engine) {
	r.GET("/chat/completions", func(c *gin.Context) {
		var req openai.ChatCompletionRequest
		err := c.BindJSON(&req)
		if err != nil {
			// TODO: log
			return
		}
		switch req.Model {
		case ChatGLM.ModelChatGLMPro, ChatGLM.ModelChatGLMStd,
			ChatGLM.ModelChatGLMLite, ChatGLM.ModelTurbo:
			ChatGLM_ChatCompletionHandler(c, req)
		}
	})
}
