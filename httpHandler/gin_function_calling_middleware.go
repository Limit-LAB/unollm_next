package httpHandler

import (
	"encoding/json"
	"io"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/utils"
)

func FunctionCallingMiddleware(c *gin.Context, req openai.ChatCompletionRequest, m func(req openai.ChatCompletionRequest) openai.ChatCompletionResponse) bool {
	pre_fc_stream := req.Stream
	fclayer := functionCallingRequestMake(&req)
	if fclayer {
		res := m(req)
		functionCallingResponseHande(c, pre_fc_stream, &req, res)
		return true
	}
	return false
}

func functionCallingRequestMake(req *openai.ChatCompletionRequest) bool {
	if len(req.Tools) != 0 && req.Messages[len(req.Messages)-1].Role == "user" {
		req.Stream = false
		userPrompt := req.Messages[len(req.Messages)-1].Content
		prefix := `you are an agent design to answer questions directly or to call other tools to answer questions,
		if you are using tools, **no need** to answer questions,
		you **only** need to answer how to call functions in **raw json format**.
		
		tools provided:

		`
		tools_json_byte, err := json.Marshal(req.Tools)
		if err != nil {
			log.Fatal(err)
		}
		tools_json_string := string(tools_json_byte)
		prefix += tools_json_string
		prefix += `
**you can ONLY use tools provided.**
if you deside to use tools, you MUST answer with following json format:
` + "```\n{\n	\"function\": {\n        \"name\": \"the name of the function\",\n        \"arguments\": \"the string of json object of the parameters you think is the best to use\"\n    }\n}\n```" +
			`if you decided to call functions, you **ONLY** need to answer raw json

the question of user is: 
`
		prefix += userPrompt
		req.Messages[len(req.Messages)-1].Content = prefix
		return true
	}
	return false
}

func functionCallingResponseHande(c *gin.Context, stream bool, req *openai.ChatCompletionRequest, resp openai.ChatCompletionResponse) {
	resp2 := resp
	var err error
	var function_calling_args string

	function_calling := getParams(`\{[ \n]*"function"[ \n]*:[ \n]*\{[ \n]*"name"[ \n]*:[ \n]*(?P<NAME>".*")[ \n]*,[ \n]*"arguments"[ \n]*:[ \n]*(?P<ARGS>".*")`, resp.Choices[0].Message.Content)
	log.Printf("%#v\n", function_calling)
	function_calling_name, ok := function_calling["NAME"]
	__content := resp.Choices[0].Message.Content
	if !ok {
		goto ERROR_HANDLING
	}
	function_calling_name, err = strconv.Unquote(function_calling_name)
	if err != nil {
		goto ERROR_HANDLING
	}
	function_calling_args, ok = function_calling["ARGS"]
	if !ok {
		goto ERROR_HANDLING
	}
	function_calling_args, err = strconv.Unquote(function_calling_args)
	if err != nil {
		goto ERROR_HANDLING
	}

	if stream {
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
	} else {
		resp2.Choices[0].Message.ToolCalls = []openai.ToolCall{{
			ID:   "unollm_simulated_" + strconv.FormatInt(time.Now().Unix(), 16),
			Type: openai.ToolType("function"),
			Function: openai.FunctionCall{
				Name:      function_calling_name,
				Arguments: function_calling_args,
			},
		}}
		resp2.Choices[0].Message.Content = ""

		resp2.Choices[0].Message.Role = openai.ChatMessageRoleAssistant
		resp2.Choices[0].Message.Content = ""
		resp2.Choices[0].FinishReason = "tool_calls"
		resp2.Choices[0].Message.ToolCalls = []openai.ToolCall{{
			ID:   "unollm_simulated_" + strconv.FormatInt(time.Now().Unix(), 16),
			Type: openai.ToolType("function"),
			Function: openai.FunctionCall{
				Name:      function_calling_name,
				Arguments: function_calling_args,
			},
		}}
		c.JSON(200, resp2)
	}
	return
ERROR_HANDLING:
	log.Println("function calling error, getting response", __content)
	c.JSON(200, resp)
}
