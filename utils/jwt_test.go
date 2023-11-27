package utils_test

import (
	"github.com/joho/godotenv"
	"go.limit.dev/unollm/provider/ChatGLM"
	"log"
	"os"
	"testing"
)

func TestJWT(t *testing.T) {
	godotenv.Load("../.env")

	zhipuaiApiKey := os.Getenv("TEST_ZHIPUAI_API")
	body := ChatGLM.ChatCompletionRequest{
		Incremental: true,
		Prompt: []ChatGLM.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "我问丁真你是哪个省的，为什么丁真回答 “我是妈妈生的？” 请给出我200字以上的答案。",
			},
		},
	}
	cli := ChatGLM.NewClient(zhipuaiApiKey)

	llm, res, err := cli.ChatCompletionStreamingRequest(body, ChatGLM.ModelTurbo)
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
