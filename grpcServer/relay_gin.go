package grpcServer

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
	"go.limit.dev/unollm/model"
	"go.limit.dev/unollm/relay"
	"go.limit.dev/unollm/relay/reqTransformer"
	"go.limit.dev/unollm/relay/respTransformer"
	"go.limit.dev/unollm/utils"
)

func getProvider(m string) (string, error) {
	if strings.Contains(m, "glm") {
		return "chatglm", nil
	}
	if strings.Contains(m, "Baichuan") {
		return "baichuan", nil
	}
	if strings.Contains(m, "gpt") {
		return "openai", nil
	}
	if strings.Contains(m, "moonshot") {
		return "moonshot", nil
	}
	return "", errors.New("could not get provider")
}

func getEmbeddingProvider(m string) (string, error) {
	if m == "embedding-2" {
		return "chatglm", nil
	}
	return "openai", nil
}

func internalServerError(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"error": err.Error(),
	})
	log.Println(err)
	c.Abort()
}

func autoErr(c *gin.Context, err error) bool {
	if err != nil {
		internalServerError(c, err)
		return true
	}
	return false
}

func RegisterRoute(r *gin.Engine) {
	v1 := r.Group("/v1")
	{
		v1.POST("/chat/completions", func(c *gin.Context) {
			var req openai.ChatCompletionRequest
			if autoErr(c, c.BindJSON(&req)) {
				return
			}
			// TODO: Model Compatitable
			provider, err := getProvider(req.Model)
			if err != nil {
				internalServerError(c, err)
				return
			}
			header := c.Request.Header["Authorization"]
			auth := header[0][7:]
			s, err := reqTransformer.ChatGPTToGrpcRequest(provider, req.Model, auth, req)
			if err != nil {
				internalServerError(c, err)
				return
			}
			mockServer := UnoForwardServer{}
			if req.Stream {
				mockServerPipe := utils.MockServerStream{
					Stream: make(chan *model.PartialLLMResponse, 1000),
				}
				mockServer.StreamRequestLLM(s, &mockServerPipe)
				respTransformer.GrpcStreamToChatGPT(c, req.Model, mockServerPipe.Stream)
				return
			} else {
				res, err := mockServer.BlockingRequestLLM(c, s)
				if err != nil {
					internalServerError(c, err)
					return
				}
				ores := respTransformer.GrpcToChatGPTCompletion(req.Model, res)
				jres, err := json.Marshal(ores)
				if err != nil {
					internalServerError(c, err)
					return
				}
				c.JSON(200, jres)
			}
		})

		v1.POST("/embeddings", func(c *gin.Context) {
			var req relay.CommonEmbdReq
			if autoErr(c, c.BindJSON(&req)) {
				return
			}
			provider, err := getEmbeddingProvider(req.Model)
			if err != nil {
				internalServerError(c, err)
				return
			}
			header := c.Request.Header["Authorization"]
			auth := header[0][7:]
			mockserver := UnoEmbeddingForwardServer{}
			res, err := mockserver.EmbeddingRequestLLM(context.Background(), &model.EmbeddingRequest{
				EmbeddingRequestInfo: &model.EmbeddingRequestInfo{
					LlmApiType: provider,
					Model:      req.Model,
					Token:      auth,
				},
				Text: req.Input,
			})
			if err != nil {
				internalServerError(c, err)
				return
			}
			oores := openai.EmbeddingResponse{
				Object: "list",
				Data: []openai.Embedding{
					{
						Object:    "embedding",
						Index:     0,
						Embedding: res.Vectors,
					},
				},
				Model: openai.EmbeddingModel(req.Model),
				Usage: openai.Usage{
					PromptTokens: int(res.Usage.PromptToken),
					TotalTokens:  int(res.Usage.TotalToken),
				},
			}
			c.JSON(200, oores)
		})
	}
}
