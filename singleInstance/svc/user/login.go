package user

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.limit.dev/unollm/singleInstance/model/apimodel"
	"go.limit.dev/unollm/singleInstance/model/dbmodel"
	"go.limit.dev/unollm/singleInstance/shared"
	"go.limit.dev/unollm/singleInstance/utils"
	"golang.org/x/crypto/bcrypt"
)

func login(c *gin.Context) {
	req, aborted := utils.GinReqJson[apimodel.UserCreatePostRequest](c)
	if aborted {
		return
	}
	var user dbmodel.User
	err := shared.GetDB().Where("`username` = ?", req.Username).First(&user).Error
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Hashed), []byte(req.Password))
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "invalid user credentials"})
		return
	}
	token := dbmodel.UserToken{
		UserId: user.ID,
		Token:  uuid.NewString(),
	}
	err = shared.GetDB().Create(&token).Error
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, apimodel.UserLoginPost200Response{Token: token.Token})
}
