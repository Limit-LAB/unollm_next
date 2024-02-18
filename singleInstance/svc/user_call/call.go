package user_call

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.limit.dev/unollm/singleInstance/model/dbmodel"
	"go.limit.dev/unollm/singleInstance/shared"
	"strings"
)

var ErrEmptyToken = errors.New("empty token")

func TokenToUser(token string) (dbmodel.User, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return dbmodel.User{}, ErrEmptyToken
	}
	var tk dbmodel.UserToken
	err := shared.GetDB().Where("`token` = ?", token).First(&tk).Error
	if err != nil {
		return dbmodel.User{}, err
	}
	var user dbmodel.User
	err = shared.GetDB().Where("`id` = ?", tk.UserId).First(&user).Error
	if err != nil {
		return dbmodel.User{}, err
	}
	return user, nil
}

const GinCtxUserKey = "user"

func MidGetUserInfoFromUserToken(c *gin.Context) {
	token := c.GetHeader("Authorization")
	user, err := TokenToUser(token)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Set(GinCtxUserKey, user)
	c.Next()
}

func GetUserInfoFromGinCtx(c *gin.Context) dbmodel.User {
	return c.MustGet(GinCtxUserKey).(dbmodel.User)
}
