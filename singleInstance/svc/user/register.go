package user

import (
	"github.com/gin-gonic/gin"
	"go.limit.dev/unollm/singleInstance/model/apimodel"
	"go.limit.dev/unollm/singleInstance/model/dbmodel"
	"go.limit.dev/unollm/singleInstance/shared"
	"go.limit.dev/unollm/singleInstance/utils"
	"golang.org/x/crypto/bcrypt"
)

func create(c *gin.Context) {
	req, abort := utils.GinReqJson[apimodel.UserCreatePostRequest](c)
	if abort {
		return
	}
	bs, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	user := dbmodel.User{
		Username: req.Username,
		Hashed:   string(bs),
	}
	err = shared.GetDB().Create(&user).Error
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}
