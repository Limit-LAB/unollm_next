package Baichuan

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func Test(t *testing.T) {
	req := BaichuanRequestBody{
		Model: "Baichuan2",
		Messages: []BaichuanMessage{
			{Role: "user", Content: "大家好啊"},
		},
		Stream:      false,
		Temperature: 0.5,
		TopP:        0.5,
		TopK:        5,
		WithSearch:  false,
	}
	godotenv.Load("../../.env")
	key := os.Getenv("TEST_BAICHUAN_API")

	c := NewClient(key)
	res, err := c.ChatCompletion(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res, err)
}
