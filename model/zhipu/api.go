package zhipu

const (
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"

	ModelChatGLMPro  = "chatglm_pro"
	ModelChatGLMStd  = "chatglm_std"
	ModelChatGLMLite = "chatglm_lite"

	ModelTurbo = "chatglm_turbo"
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRef struct {
	Enable      bool   `json:"enable"`
	SearchQuery string `json:"search_query"`
}

// ChatCompletionRequest represents a request structure for chat completion API.
type ChatCompletionRequest struct {
	Prompt      []ChatCompletionMessage `json:"prompt"`
	Temperature float32                 `json:"temperature,omitempty"`
	TopP        float32                 `json:"top_p,omitempty"`
	RequestId   string                  `json:"request_id"`
	// in SSE mode, incremental means whether to return incremental result or add to previous result
	Incremental bool              `json:"incremental"`
	ReturnType  string            `json:"return_type,omitempty"`
	Ref         ChatCompletionRef `json:"ref"`
}

// ChatCompletionResponse represents a response structure for chat completion API.
type ChatCompletionResponse struct {
	Event string                      `json:"event"`
	Data  *ChatCompletionResponseData `json:"data"`

	ErrorCode int    `json:"code"`
	ErrorMsg  string `json:"msg"`
	Success   bool   `json:"success"`
}

type ChatCompletionResponseData struct {
	TaskId  string                 `json:"task_id"`
	Usage   Usage                  `json:"usage"`
	Choices []ChatCompletionChoice `json:"choices"`
}

type ChatCompletionChoice struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
