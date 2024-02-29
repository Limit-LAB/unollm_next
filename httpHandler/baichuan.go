package httpHandler

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/relay/reqTransformer"
	"go.limit.dev/unollm/relay/respTransformer"
)

func getParams(regEx, url string) (paramsMap map[string]string) {
	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}

func Baichuan_ChatCompletionHandler(c *gin.Context, tx KeyTransformer, req openai.ChatCompletionRequest) {
	cli := NewBaichuanClient(c, tx)
	if cli == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "no api key provided",
		})
		c.Abort()
		return
	}

	// addtion layer to support function calling
	fc := FunctionCallingMiddleware(c, req, func(req openai.ChatCompletionRequest) openai.ChatCompletionResponse {
		rst, err := cli.ChatCompletion(reqTransformer.BaiChuanFromOpenAIChatCompletionReq(req))
		if err != nil {
			internalServerError(c, err)
			return openai.ChatCompletionResponse{}
		}
		return respTransformer.BaiChuanToOpenAICompletion(rst)
	})
	if fc {
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
