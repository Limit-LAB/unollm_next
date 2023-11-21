package utils

import (
	"log"
	"os"
	"testing"

	"limit.dev/unollm/model/zhipu"

	"github.com/joho/godotenv"
)

func TestJWT(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	zhipuaiApiKey := os.Getenv("TEST_ZHIPUAI_API")
	body := zhipu.ChatCompletionRequest{
		Incremental: true,
		Prompt: []zhipu.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "我问丁真你是哪个省的，为什么丁真回答 “我是妈妈生的？” 请给出我200字以上的答案。",
			},
		},
	}

	llm, res, err := GLMStreamingRequest(body, "chatglm_turbo", zhipuaiApiKey)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case llm_message := <-llm:
			log.Print(llm_message)
		case result := <-res:
			log.Print(result)
			goto BYE
		}
	}

BYE:
	return
}
