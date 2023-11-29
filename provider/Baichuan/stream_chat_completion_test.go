package Baichuan

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestStream(t *testing.T) {
	req := BaichuanRequestBody{
		Model: "Baichuan2",
		Messages: []BaichuanMessage{
			{Role: "user", Content: "明天下大雨，我该不该打伞？"},
		},
		Stream:      true,
		Temperature: 0.5,
		TopP:        0.5,
		TopK:        5,
		WithSearch:  false,
	}
	godotenv.Load("../../.env")
	key := os.Getenv("TEST_BAICHUAN_API")

	c := NewClient(key)
	res, err := c.ChatCompletionStreamingRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	for r := range res {
		if r.Choices[0].Finish_reason == "stop" {
			break
		}
		t.Log(r)
	}
}
