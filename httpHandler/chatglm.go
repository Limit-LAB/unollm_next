package httpHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/provider/ChatGLM"
	"go.limit.dev/unollm/relay/reqTransformer"
	"go.limit.dev/unollm/relay/respTransformer"
)

func ChatGLM_ChatCompletionHandler(c *gin.Context, tx KeyTransformer,req openai.ChatCompletionRequest) {
	cli := NewChatGLMClient(c, tx)
	if cli == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "no api key provided",
		})
		c.Abort()
		return
	}

	zpReq := reqTransformer.ChatGLMFromOpenAIChatCompletionReq(req)
	if req.Stream {
		resp, err := cli.ChatCompletionStreamingRequest(zpReq)

		if err != nil {
			internalServerError(c, err)
			return
		}
		respTransformer.ChatGLMToOpenAIStream(c, resp)
		return
	}

	rst, err := cli.ChatCompletion(zpReq)
	if err != nil {
		internalServerError(c, err)
		return
	}
	c.JSON(200, respTransformer.ChatGLMToOpenAICompletion(rst))
}

func ChatGLM_EmbeddingHandler(c *gin.Context, req OpenAIEmbeddingRequest) {
	cli := NewChatGLMClient(c)
	res, err := cli.EmbeddingRequest(
		ChatGLM.EmbeddingRequest{
			Input: req.Input,
			Model: req.Model[9:],
		},
	)
	if err != nil {
		internalServerError(c, err)
		return
	}
	c.JSON(200, res)
}
