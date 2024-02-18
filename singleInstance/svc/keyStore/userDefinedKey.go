package keyStore

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.limit.dev/unollm/singleInstance/model/apimodel"
	"go.limit.dev/unollm/singleInstance/model/dbmodel"
	"go.limit.dev/unollm/singleInstance/shared"
	"go.limit.dev/unollm/singleInstance/utils"
)

func (svc *KeyStoreSvc) newApi(c *gin.Context) {
	req, aborted := utils.GinReqJson[apimodel.KeyStoreNewApiPostRequest](c)
	if aborted {
		return
	}
	key := req.GetKey()
	if key == "" {
		key = uuid.NewString()
	}
	udk := dbmodel.UserDefinedKey{
		Owner: getUidFrom(c),
		Key:   key,
	}
	err := shared.GetDB().Save(&udk).Error
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"key": udk.Key,
		"id":  int32(udk.ID),
	})

}

func (svc *KeyStoreSvc) listUserDefinedKeys(c *gin.Context) {
	uid := getUidFrom(c)
	var keys []dbmodel.UserDefinedKey
	err := shared.GetDB().Where("owner = ?", uid).Find(&keys).Error
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	var res []apimodel.KeyStoreUserDefinedKeysGet200ResponseInner
	for _, key := range keys {
		res = append(res, apimodel.KeyStoreUserDefinedKeysGet200ResponseInner{
			Key: key.Key,
		})
	}
	c.JSON(200, res)
}
