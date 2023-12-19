package relay

import (
	"go.limit.dev/unollm/provider/Baichuan"
)

func BaiChuanChatCompletionRequest(cli *Baichuan.Client, req Baichuan.ChatCompletionRequest) (Baichuan.ChatCompletionResponse, error) {
	req.Stream = false
	return cli.ChatCompletion(req)
}

func BaiChuanChatCompletionStreamingRequest(cli *Baichuan.Client, req Baichuan.ChatCompletionRequest) (chan Baichuan.StreamResponse, error) {
	req.Stream = true
	return cli.ChatCompletionStreamingRequest(req)
}
