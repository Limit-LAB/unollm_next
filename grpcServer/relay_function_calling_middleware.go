package grpcServer

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"

	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
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

func functionCallingRequestMake(req *model.LLMRequestSchema) bool {
	if len(req.LlmRequestInfo.Functions) != 0 && req.Messages[len(req.Messages)-1].Role == "user" {
		userPrompt := req.Messages[len(req.Messages)-1].Content
		prefix := `you are an agent design to answer questions directly or to call other tools to answer questions,
		if you are using tools, **no need** to answer questions,
		you **only** need to answer how to call functions in **raw json format**.
		
		tools provided:

		`
		info := req.LlmRequestInfo
		tools := make([]openai.Tool, len(info.Functions))
		for i, f := range info.Functions {
			tools[i] = openai.Tool{
				Type: "function",
				Function: openai.FunctionDefinition{
					Name:        f.Name,
					Description: f.Description,
					Parameters: map[string]any{
						"type":       "object",
						"properties": map[string]any{},
						"required":   f.Requireds,
					},
				},
			}
			for j := 0; j < len(f.Parameters); j++ {
				tools[i].Function.Parameters.(map[string]any)["properties"].(map[string]any)[f.Parameters[j].Name] = map[string]any{
					"type":        f.Parameters[j].Type,
					"description": f.Parameters[j].Description,
				}
				if len(f.Parameters[j].Enums) != 0 {
					tools[i].Function.Parameters.(map[string]any)["properties"].(map[string]any)[f.Parameters[j].Name].(map[string]any)["enum"] = f.Parameters[j].Enums
				}
			}
			log.Printf("%#v\n", tools[i])
		}

		tools_json_byte, err := json.Marshal(tools)
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

func functionCallingResponseHandle(resp *model.LLMResponseSchema) {
	var err error
	var function_calling_args string

	function_calling := getParams(`\{[ \n]*"function"[ \n]*:[ \n]*\{[ \n]*"name"[ \n]*:[ \n]*(?P<NAME>".*")[ \n]*,[ \n]*"arguments"[ \n]*:[ \n]*(?P<ARGS>".*")`, resp.Message.Content)
	log.Printf("%#v\n", function_calling)
	function_calling_name, ok := function_calling["NAME"]
	__content := resp.Message.Content
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
	resp.Message.Content = ""
	resp.ToolCalls = []*model.ToolCall{
		{
			Name:      function_calling_name,
			Arguments: function_calling_args,
		},
	}
	return
ERROR_HANDLING:
	log.Println("function calling error, getting response", __content)
}

func functionCallingResponseToStream(resp *model.LLMResponseSchema, sv model.UnoLLMv1_StreamRequestLLMServer) {
	sv.Send(&model.PartialLLMResponse{
		ToolCalls: resp.ToolCalls,
		Response:  &model.PartialLLMResponse_Content{Content: resp.Message.Content},
	})
	sv.Send(&model.PartialLLMResponse{
		Response:      &model.PartialLLMResponse_Done{},
		LlmTokenCount: resp.LlmTokenCount,
	})
	return
}
