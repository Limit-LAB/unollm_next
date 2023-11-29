package Baichuan

type BaichuanMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type BaichuanRequestBody struct {
	Model       string            `json:"model"`
	Messages    []BaichuanMessage `json:"messages"`
	Stream      bool              `json:"stream,omitempty"`
	Temperature float32           `json:"temperature,omitempty"`
	TopP        float32           `json:"top_p,omitempty"`
	TopK        int               `json:"top_k,omitempty"`
	WithSearch  bool              `json:"with_search_enhance,omitempty"`
}

/**/

type BaichuanBlockingResponseChoices struct {
	Finish_reason string            `json:"finish_reason"`
	Index         int               `json:"index"`
	Message       []BaichuanMessage `json:"message"`
}

type BaichuanStreamResponseChoices struct {
	Finish_reason string          `json:"finish_reason"`
	Index         int             `json:"index"`
	Delta         BaichuanMessage `json:"delta"`
}

type BaichuanResponseUsage struct {
	Completion_tokens int `json:"completion_tokens"`
	Prompt_tokens     int `json:"prompt_tokens"`
	Total_tokens      int `json:"total_tokens"`
}

type BaichuanBlockingResponseBody struct {
	Id      string                            `json:"id"`
	Created int                               `json:"created"`
	Choices []BaichuanBlockingResponseChoices `json:"choices"`
	Model   string                            `json:"model"`
	Object  string                            `json:"object"`
	Usage   BaichuanResponseUsage             `json:"usage"`
}

type BaichuanStreamResponseBody struct {
	Id      string                          `json:"id"`
	Created int                             `json:"created"`
	Choices []BaichuanStreamResponseChoices `json:"choices"`
	Model   string                          `json:"model"`
	Object  string                          `json:"object"`
}
