package httpHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/relay/reqTransformer"
	"go.limit.dev/unollm/relay/respTransformer"
)

func Baichuan_ChatCompletionHandler(c *gin.Context, tx KeyTransformer, req openai.ChatCompletionRequest) {
	cli := NewBaichuanClient(c, tx)
	if cli == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "no api key provided",
		})
		c.Abort()
		return
	}

	zpReq := reqTransformer.BaiChuanFromOpenAIChatCompletionReq(req)
	if req.Stream {
		resp, err := cli.ChatCompletionStreamingRequest(zpReq)

		if err != nil {
			internalServerError(c, err)
			return
		}
		respTransformer.BaiChuanToOpenAIStream(c, resp)
		return
	}

	rst, err := cli.ChatCompletion(zpReq)
	if err != nil {
		internalServerError(c, err)
		return
	}
	c.JSON(200, respTransformer.BaiChuanToOpenAICompletion(rst))
}
