package utils_test

import (
	"github.com/joho/godotenv"
	"limit.dev/unollm/utils"
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

	llm, res, err := utils.GLMStreamingRequest(body, zhipu.ModelTurbo, zhipuaiApiKey)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case llmMessage := <-llm:
			log.Print(llmMessage)
		case result := <-res:
			log.Print(result)
			goto BYE
		}
	}

BYE:
	return
}
