package grpcServer

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/provider/ChatGLM"
)

type Callings struct {
}

func removeWhiteSpaces(str string) string {
	var b strings.Builder
	b.Grow(len(str))
	for _, ch := range str {
		if ch != ' ' && ch != '\t' && ch != '\n' {
			b.WriteRune(ch)
		}
	}
	return b.String()
}
func getParams(url string) (params []ChatGLM.GLMFunctionCall, err error) {
	url = removeWhiteSpaces(url)
	err = json.Unmarshal([]byte(url), &params)
	if err != nil {
		// ```json<>```
		uurl1, found1 := strings.CutPrefix(url, "```json")
		if found1 {
			uurl1, _ = strings.CutSuffix(uurl1, "```")
			err = json.Unmarshal([]byte(uurl1), &params)
			return
		}
		uurl2, found2 := strings.CutPrefix(url, "```")
		if found2 {
			uurl2, _ = strings.CutSuffix(uurl2, "```")
			err = json.Unmarshal([]byte(uurl2), &params)
			return
		}
	}
	return
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
				Function: &openai.FunctionDefinition{
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
			log.Printf("find function: %#v", tools[i].Function)
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
` + "```[\n{\n        \"name\": \"the name of the function\",\n        \"arguments\": \"the string of json object of the parameters you think is the best to use\"\n    }\n, /* function 2, function 3, etc...*/ ]```" +
			`if you decided to call functions, you **ONLY** need to answer raw json
`
		prefix += userPrompt
		req.Messages[len(req.Messages)-1].Content = prefix
		return true
	}
	return false
}

func functionCallingResponseHandle(resp *model.LLMResponseSchema) {
	function_calling, err := getParams(resp.Message.Content)
	if err != nil {
		log.Fatal(err)
	}
	resp.Message.Content = ""
	resp.ToolCalls = []*model.ToolCall{}
	for _, f := range function_calling {
		resp.ToolCalls = append(resp.ToolCalls, &model.ToolCall{
			Id:        fmt.Sprintf("unollm_adapter_%d", time.Now().Unix()),
			Name:      f.Name,
			Arguments: f.Arguments,
		})
	}
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
}
