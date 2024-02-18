package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GinReqJson[T any](c *gin.Context) (val T, aborted bool) {
	if c.IsAborted() {
		return
	}
	err := c.ShouldBindJSON(&val)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return val, true
	}
	return
}
