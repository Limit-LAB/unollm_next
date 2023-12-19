package Baichuan

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream,omitempty"`
	Temperature float32   `json:"temperature,omitempty"`
	TopP        float32   `json:"top_p,omitempty"`
	TopK        int       `json:"top_k,omitempty"`
	WithSearch  bool      `json:"with_search_enhance,omitempty"`
}

type Choice struct {
	Finish_reason string  `json:"finish_reason"`
	Index         int     `json:"index"`
	Message       Message `json:"message"`
}

type StreamChoice struct {
	Finish_reason string  `json:"finish_reason"`
	Index         int     `json:"index"`
	Delta         Message `json:"delta"`
}

type Usage struct {
	Completion_tokens int `json:"completion_tokens"`
	Prompt_tokens     int `json:"prompt_tokens"`
	Total_tokens      int `json:"total_tokens"`
}

type ChatCompletionResponse struct {
	Id      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type StreamResponse struct {
	Id      string         `json:"id"`
	Created int            `json:"created"`
	Choices []StreamChoice `json:"choices"`
	Model   string         `json:"model"`
	Object  string         `json:"object"`
}
