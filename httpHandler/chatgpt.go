package httpHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/relay"
	"go.limit.dev/unollm/relay/respTransformer"
	"log"
)

func ChatGPT_ChatCompletitionsHandler(c *gin.Context, tx KeyTransformer, req openai.ChatCompletionRequest) {
	cli := NewOpenAIClient(c, tx)
	if cli == nil {
		c.JSON(401, gin.H{
			"error": "no api key provided",
		})
		c.Abort()
		return
	}
	if !req.Stream {
		rsp, err := relay.OpenAIChatCompletionRequest(cli, req)
		if err != nil {
			internalServerError(c, err)
			return
		}
		c.JSON(200, rsp)
		return
	}

	rsp, err := relay.OpenAIChatCompletionStream(cli, req)
	if err != nil {
		internalServerError(c, err)
		return
	}
	respTransformer.ChatGPTToOpenAIChatCompletionStream(c, rsp)
}

func ChatGPT_CompletitionsHandler(c *gin.Context, tx KeyTransformer, req openai.CompletionRequest) {
	cli := NewOpenAIClient(c, tx)
	if cli == nil {
		c.JSON(401, gin.H{
			"error": "no api key provided",
		})
		c.Abort()
		return
	}
	if !req.Stream {
		rsp, err := relay.OpenAICompletionRequest(cli, req)
		if err != nil {
			internalServerError(c, err)
			return
		}
		c.JSON(200, rsp)
		return
	}

	rsp, err := relay.OpenAICompletionStream(cli, req)
	if err != nil {
		internalServerError(c, err)
		return
	}
	respTransformer.ChatGPTToOpenAICompletionStream(c, rsp)
}

func internalServerError(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"error": err.Error(),
	})
	log.Println(err)
	c.Abort()
	return
}
