package httpHandler

import (
	"encoding/json"
	"log"

	"github.com/sashabaranov/go-openai"
)

func FunctionCallingRequestMake(req *openai.ChatCompletionRequest) bool {
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
