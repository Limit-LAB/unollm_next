package respTransformer

import (
	"encoding/json"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/provider/ChatGLM"
	"go.limit.dev/unollm/utils"
)

func chatGlmChoicesToOpenAIChoices(choices []ChatGLM.ChatCompletionChoice) []openai.ChatCompletionChoice {
	var openAIChoices []openai.ChatCompletionChoice
	for _, choice := range choices {
		openAIChoices = append(openAIChoices, openai.ChatCompletionChoice{
			Index: choice.Index,
			Message: openai.ChatCompletionMessage{
				Role:    choice.Message.Role,
				Content: choice.Message.Content,
			},
		})
	}
	return openAIChoices

}

func chatGlmDeltaChoicesToOpenAIChoices(choices []ChatGLM.ChatCompletionStreamingChoice) []openai.ChatCompletionStreamChoice {
	var openAIChoices []openai.ChatCompletionStreamChoice
	for _, choice := range choices {
		toolCalls := make([]openai.ToolCall, len(choice.Delta.ToolCalls))
		for i, f := range choice.Delta.ToolCalls {
			toolCalls[i].ID = f.Id
			toolCalls[i].Type = openai.ToolType(f.Type)
			toolCalls[i].Function.Arguments = f.Function.Arguments
			toolCalls[i].Function.Name = f.Function.Name
		}

		openAIChoices = append(openAIChoices, openai.ChatCompletionStreamChoice{
			Index: choice.Index,
			Delta: openai.ChatCompletionStreamChoiceDelta{
				Content:   choice.Delta.Content,
				Role:      choice.Delta.Role,
				ToolCalls: toolCalls,
			},
			FinishReason: openai.FinishReason(choice.FinishReason),
		})
	}
	return openAIChoices

}

func ChatGLMToOpenAICompletion(res ChatGLM.ChatCompletionResponse) openai.ChatCompletionResponse {
	return openai.ChatCompletionResponse{
		ID:      res.Id,
		Choices: chatGlmChoicesToOpenAIChoices(res.Choices),
		Usage: openai.Usage{
			PromptTokens:     res.Usage.PromptTokens,
			TotalTokens:      res.Usage.TotalTokens,
			CompletionTokens: res.Usage.CompletionTokens,
		},
	}
}

func ChatGLMToOpenAIStream(c *gin.Context, _r *ChatGLM.ChatCompletionStreamingResponse) {
	llm, result := _r.ResponseCh, _r.FinishCh
	utils.SetEventStreamHeaders(c)
	c.Stream(func(w io.Writer) bool {
		select {
		case data := <-llm:
			response := openai.ChatCompletionStreamResponse{
				Object:  "chat.completion.chunk",
				Model:   data.Model,
				Created: data.Created,
				Choices: chatGlmDeltaChoicesToOpenAIChoices(data.Choices),
			}
			jsonResponse, _ := json.Marshal(response)
			c.Render(-1, utils.CustomEvent{Data: "data: " + string(jsonResponse)})
			return true
		case _ = <-result:
			c.Render(-1, utils.CustomEvent{Data: "data: [DONE]"})
			return false
		}
	})
}
