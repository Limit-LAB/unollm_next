package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"testing"

	"limit.dev/unollm/model/zhipu"
)

func TestJWT(t *testing.T) {
	godotenv.Load("../.env")

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
