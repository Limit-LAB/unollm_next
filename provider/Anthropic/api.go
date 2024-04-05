package anthropic

import (
	"encoding/xml"
	"net/http"
)

type Client struct {
	key string
	hc  *http.Client
}

func NewClient(key string) *Client {
	return &Client{
		key: key,
		hc:  &http.Client{},
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"` // 4096
	Stream      bool      `json:"stream,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
	TopK        float64   `json:"top_k,omitempty"`
}

type Content struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}
type ChatCompletionResponse struct {
	Content      []Content `json:"content"`
	Id           string    `json:"id"`
	Model        string    `json:"model"`
	Role         string    `json:"role"`
	StopReason   string    `json:"stop_reason"`
	StopSequence string    `json:"stop_sequence,omitempty"`
	Type         string    `json:"type"`
	Usage        Usage     `json:"usage"`
}

const OPUS = "claude-3-opus-20240229"
const SONNET = "claude-3-sonnet-20240229"

type FunctionCalling struct {
	XMLName xml.Name `xml:"config"` // 指定最外层的标签为config
}

// type StreamMessageStartMessage struct {
// 	Id    string `json:"id"`
// 	Type  string `json:"type"`
// 	Role  string `json:"role"`
// 	Model string `json:"model"`
// 	Usage Usage  `json:"usage"`
// }

// type StreamMessageStart struct {
// 	Type    string                    `json:"type"`
// 	Message StreamMessageStartMessage `json:"message"`
// }

type StreamResponseDelta struct {
	Type         string `json:"type"`
	Text         string `json:"text"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Usage        Usage  `json:"usage"`
}

type StreamResponse struct {
	Type  string              `json:"type"`
	Index int                 `json:"index"`
	Delta StreamResponseDelta `json:"delta"`
}
