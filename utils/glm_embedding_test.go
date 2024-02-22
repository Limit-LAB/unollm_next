package utils_test

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"go.limit.dev/unollm/provider/ChatGLM"
)

func TestEmbedding(t *testing.T) {
	godotenv.Load("../.env")

	zhipuaiApiKey := os.Getenv("TEST_ZHIPUAI_API")
	cli := ChatGLM.NewClient(zhipuaiApiKey)
	body := ChatGLM.EmbeddingRequest{
		Input: "丁真的家乡来自四川理塘",
		Model: "embedding-2",
	}
	result, err := cli.EmbeddingRequest(body)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(result)
}
