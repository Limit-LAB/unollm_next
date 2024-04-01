package ChatGLM

import (
	"go.limit.dev/unollm/model"
)

const (
	ModelGLM3Turbo = "glm-3-turbo"
	ModelGLM4      = "glm-4"
	ModelGLM4V     = "glm-4v"
)

type ChatCompletionTextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Image struct {
	Url string `json:"url"`
}
type ChatCompletionImageContent struct {
	Type     string `json:"type"`
	ImageUrl Image  `json:"image_url"`
}

type ChatCompletionMessage struct {
	Role      string        `json:"role"`
	Content   any           `json:"content"`
	ToolCalls []GLMToolCall `json:"tool_calls,omitempty"`
}

type GLMToolCall struct {
	Id       string          `json:"id"`
	Type     string          `json:"type"`
	Function GLMFunctionCall `json:"function"`
}

type GLMFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatCompletionRef struct {
	Enable      bool   `json:"enable"`
	SearchQuery string `json:"search_query"`
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	Model       string                  `json:"model"`
	Messages    []ChatCompletionMessage `json:"messages"`
	DoSample    bool                    `json:"do_sample,omitempty"`
	Temperature float32                 `json:"temperature,omitempty"`
	TopP        float32                 `json:"top_p,omitempty"`
	MaxTokens   int                     `json:"max_tokens,omitempty"`
	RequestId   string                  `json:"request_id,omitempty"`
	Stream      bool                    `json:"stream,omitempty"`
	Stop        []string                `json:"stop,omitempty"`
	Tools       []GLMTool               `json:"tools,omitempty"`
	ToolChoice  string                  `json:"tool_choice,omitempty"`
}

type GLMTool struct {
	// only "function" "retrieval" "web_search"
	Type      string       `json:"type"`
	Function  GLMFunction  `json:"function,omitempty"`
	Retrieval any          `json:"retrieval,omitempty"`
	WebSearch GLMWebSearch `json:"web_search,omitempty"`
}

type GLMFunction struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  any    `json:"parameters"`
}
type GLMWebSearch struct {
	Enable      bool   `json:"enable"`
	SearchQuery string `json:"search_query"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	Id      string                 `json:"id"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model,omitempty"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   Usage                  `json:"usage"`
}

type ChatCompletionStreamResponse struct {
	Id      string                          `json:"id"`
	Created int64                           `json:"created"`
	Model   string                          `json:"model"`
	Choices []ChatCompletionStreamingChoice `json:"choices"`
	Usage   Usage                           `json:"usage,omitempty"`
}

type ChatCompletionStreamingChoice struct {
	Index        int                   `json:"index"`
	Delta        ChatCompletionMessage `json:"delta"`
	FinishReason FinishReason          `json:"finish_reason,omitempty"`
}

type ChatCompletionResponseData struct {
	TaskId  string                 `json:"task_id"`
	Usage   Usage                  `json:"usage"`
	Choices []ChatCompletionChoice `json:"choices"`
}

type ChatCompletionChoice struct {
	FinishReason string                `json:"finish_reason,omitempty"`
	Index        int                   `json:"index"`
	Message      ChatCompletionMessage `json:"message"`
}

type FinishReason string

const (
	FinishReasonMaxTokens FinishReason = "max_tokens"
	FinishReasonNone      FinishReason = ""
	FinishReasonStop      FinishReason = "stop"
	FinishReasonLength    FinishReason = "length"
	FinishReasonToolCalls FinishReason = "tool_calls"
)

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func (u Usage) ToGrpc() model.LLMTokenCount {
	return model.LLMTokenCount{
		TotalToken:      int64(u.TotalTokens),
		PromptToken:     int64(u.PromptTokens),
		CompletionToken: int64(u.CompletionTokens),
	}
}

type EmbeddingRequest struct {
	Input string `json:"input"`
	Model string `json:"model"`
}

type EmbeddingResponseData struct {
	Index     int       `json:"index"`
	Object    string    `json:"object"`
	Embedding []float32 `json:"embedding"`
}

type EmbeddingResponse struct {
	Model  string                  `json:"model"`
	Data   []EmbeddingResponseData `json:"data"`
	Object string                  `json:"object"`
	Usage  Usage                   `json:"usage"`
}
