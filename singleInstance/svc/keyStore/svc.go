package keyStore

import (
	"github.com/gin-gonic/gin"
	"go.limit.dev/unollm/httpHandler"
)

type KeyStoreSvc struct{}

func KeyTransformer(key string, provider string) (httpHandler.KeyTransformerResult, error) {
	return httpHandler.KeyTransformerResult{
		Key:      key,
		EndPoint: "",
	}, nil
}

func (svc *KeyStoreSvc) RegisterRouter(g gin.IRouter, prefix string) {
	_ = g.Group(prefix)

}
