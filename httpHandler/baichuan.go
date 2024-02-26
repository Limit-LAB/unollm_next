package httpHandler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/provider/Baichuan"
	"go.limit.dev/unollm/relay/reqTransformer"
	"go.limit.dev/unollm/relay/respTransformer"
	"go.limit.dev/unollm/utils"
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

// TODO: migrate common part into relay
func functionCallingMiddleware(
	req openai.ChatCompletionRequest,
	cli *Baichuan.Client,
	c *gin.Context,
) (bool, Baichuan.ChatCompletionRequest) {
	pre_fc_stream := req.Stream
	fclayer := FunctionCallingRequestMake(&req)

	zpReq := reqTransformer.BaiChuanFromOpenAIChatCompletionReq(req)

	if fclayer {
		rst, err := cli.ChatCompletion(zpReq)
		if err != nil {
			internalServerError(c, err)
			return true, zpReq
		}
		log.Printf("%#v", rst.Choices[0])
		function_calling := getParams(`\{[ \n]*"function"[ \n]*:[ \n]*\{[ \n]*"name"[ \n]*:[ \n]*(?P<NAME>".*")[ \n]*,[ \n]*"arguments"[ \n]*:[ \n]*(?P<ARGS>".*")`, rst.Choices[0].Message.Content)
		log.Printf("%#v\n", function_calling)
		function_calling_name, ok := function_calling["NAME"]
		if !ok {
			// do something
		}
		function_calling_name, err = strconv.Unquote(function_calling_name)
		if err != nil {
			// do something
		}
		function_calling_args, ok := function_calling["ARGS"]
		if !ok {
			// do something
		}
		function_calling_args, err = strconv.Unquote(function_calling_args)
		if err != nil {
			// do something
		}
		if pre_fc_stream {
			c.Stream(func(w io.Writer) bool {

				response := openai.ChatCompletionStreamResponse{
					Object:  "chat.completion.chunk",
					Model:   req.Model,
					Created: time.Now().Unix(),
					Choices: []openai.ChatCompletionStreamChoice{
						{
							Delta: openai.ChatCompletionStreamChoiceDelta{
								Content: "",
								Role:    openai.ChatMessageRoleAssistant,
								ToolCalls: []openai.ToolCall{{
									ID:   "unollm_simulated_" + strconv.FormatInt(time.Now().Unix(), 16),
									Type: openai.ToolType("function"),
									Function: openai.FunctionCall{
										Name:      function_calling_name,
										Arguments: function_calling_args,
									},
								}},
							},
						},
					},
				}
				jsonResponse, _ := json.Marshal(response)
				c.Render(-1, utils.CustomEvent{Data: "data: " + string(jsonResponse)})
				c.Render(-1, utils.CustomEvent{Data: "data: [DONE]"})
				return false
			})
			return true, zpReq
		}
		res := respTransformer.BaiChuanToOpenAICompletion(rst)
		res.Choices[0].Message.Role = openai.ChatMessageRoleAssistant
		res.Choices[0].Message.Content = ""
		res.Choices[0].FinishReason = "tool_calls"
		res.Choices[0].Message.ToolCalls = []openai.ToolCall{{
			ID:   "unollm_simulated_" + strconv.FormatInt(time.Now().Unix(), 16),
			Type: openai.ToolType("function"),
			Function: openai.FunctionCall{
				Name:      function_calling_name,
				Arguments: function_calling_args,
			},
		}}
		c.JSON(200, res)
		return true, zpReq
	}
	return false, zpReq
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
	fc, zpReq := functionCallingMiddleware(req, cli, c)
	if fc {
		return
	}

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
