package utils_test

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.limit.dev/unollm/provider/ChatGLM"
)

func TestJWT(t *testing.T) {
	godotenv.Load("../.env")

	zhipuaiApiKey := os.Getenv("TEST_ZHIPUAI_API")
	body := ChatGLM.ChatCompletionRequest{
		Model: ChatGLM.ModelGLM3Turbo,
		Messages: []ChatGLM.ChatCompletionMessage{
			{
				Role:    "user",
				Content: "我问丁真你是哪个省的，为什么丁真回答 “我是妈妈生的？” 请给出我200字以上的答案。",
			},
		},
	}
	cli := ChatGLM.NewClient(zhipuaiApiKey)
	res, err := cli.ChatCompletionStreamingRequest(body)
	if err != nil {
		log.Fatal(err)
		return
	}
	for {
		select {
		case chunk := <-res.ResponseCh:
			log.Print(chunk.Choices[0].Delta.Content)
		case result := <-res.FinishCh:
			log.Print(result)
			goto END
		}
	}
END:
	return

	// llm, res := _r.OnRecvData, _r.OnFinish

	// for {
	// 	select {
	// 	case llmMessage := <-llm:
	// 		log.Print(llmMessage)
	// 	case result := <-res:
	// 		log.Print(result)
	// 		goto BYE
	// 	}
	// }

	// BYE:
	// return
}
