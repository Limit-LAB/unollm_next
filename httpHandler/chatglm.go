package httpHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/relay/reqTransformer"
	"go.limit.dev/unollm/relay/respTransformer"
	"net/http"
)

func ChatGLM_ChatCompletionHandler(c *gin.Context, req openai.ChatCompletionRequest) {
	cli := NewChatGLMClient(c)
	if cli == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "no api key provided",
		})
		c.Abort()
		return
	}

	zpReq := reqTransformer.ChatGLMFromOpenAIChatCompletionReq(req)
	if req.Stream {
		resp, err := cli.ChatCompletionStreamingRequest(zpReq, req.Model)

		if err != nil {
			internalServerError(c, err)
			return
		}
		respTransformer.ChatGLMToOpenAIStream(c, resp)
		return
	}

	rst, err := cli.ChatCompletion(zpReq, req.Model)
	if err != nil {
		internalServerError(c, err)
		return
	}
	c.JSON(200, respTransformer.ChatGLMToOpenAICompletion(rst))
	return

}
