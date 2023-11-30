package keyStore

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type KeyStoreSvc struct{}

var ErrNoKey = errors.New("no key provided")

func (svc *KeyStoreSvc) RegisterRouter(g gin.IRouter, prefix string) {
	_ = g.Group(prefix)

}
