package utils

import (
	"limit.dev/unollm/model/zhipu"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestJWT(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	zhipuaiApiKey := os.Getenv("TEST_ZHIPUAI_API")
	body := zhipu.ChatCompletionRequest{
		Prompt: []zhipu.ChatCompletionMessage{
			zhipu.ChatCompletionMessage{
				Role:    "user",
				Content: "我问丁真你是哪个省的，为什么丁真回答 “我是妈妈生的？” 请给出我200字以上的答案。",
			},
		},
	}

	result, err := GLMBlockingRequest(body, "chatglm_turbo", zhipuaiApiKey)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", result)
}
