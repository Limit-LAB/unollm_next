package keyStore

import (
	"github.com/gin-gonic/gin"
	"go.limit.dev/unollm/httpHandler"
	"go.limit.dev/unollm/singleInstance/model/apimodel"
	"go.limit.dev/unollm/singleInstance/model/dbmodel"
	"go.limit.dev/unollm/singleInstance/shared"
	"go.limit.dev/unollm/singleInstance/utils"
	"math/rand"
)

func KeyTransformer(keyIn string, provider string) (httpHandler.KeyTransformerResult, error) {
	var userKey dbmodel.UserDefinedKey
	err := shared.GetDB().
		Where("`key` = ?", keyIn).
		First(&userKey).Error
	if err != nil {
		return httpHandler.KeyTransformerResult{}, err
	}

	var keyRelation []dbmodel.MapOriginRelation
	err = shared.GetDB().
		Where("map_id = ?", userKey.ID).
		Find(&keyRelation).Error
	if err != nil {
		return httpHandler.KeyTransformerResult{}, err
	}
	if len(keyRelation) == 0 {
		return httpHandler.KeyTransformerResult{}, ErrNoKey
	}

	var originKeyIds []uint
	for _, relation := range keyRelation {
		originKeyIds = append(originKeyIds, relation.OriginID)
	}

	var keys []dbmodel.OriginKey
	err = shared.GetDB().
		Where("id in ?", originKeyIds).
		Where("provider = ?", provider).
		Where("disabled = ?", false).
		Find(&keys).Error
	if err != nil {
		return httpHandler.KeyTransformerResult{}, err
	}
	if len(keys) == 0 {
		return httpHandler.KeyTransformerResult{}, ErrNoKey
	}

	key := keys[rand.Intn(len(keys))]

	return httpHandler.KeyTransformerResult{
		Key:      key.Key,
		EndPoint: key.EndPoint,
	}, nil
}

func (svc *KeyStoreSvc) testTransformer(c *gin.Context) {
	req, aborted := utils.GinReqJson[apimodel.KeyStoreTestTransformerPostRequest](c)
	if aborted {
		return
	}
	rst, err := KeyTransformer(req.Key, req.Provider)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	resp := apimodel.NewKeyStoreTestTransformerPost200Response(rst.Key)
	if rst.EndPoint != "" {
		resp.SetEndpoint(rst.EndPoint)
	}
	c.JSON(200, resp)
}
