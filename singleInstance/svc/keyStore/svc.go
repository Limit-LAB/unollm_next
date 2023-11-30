package keyStore

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.limit.dev/unollm/singleInstance/model/apimodel"
	"go.limit.dev/unollm/singleInstance/model/dbmodel"
	"go.limit.dev/unollm/singleInstance/shared"
	"go.limit.dev/unollm/singleInstance/utils"
)

type KeyStoreSvc struct{}

var ErrNoKey = errors.New("no key provided")

func (svc *KeyStoreSvc) RegisterRouter(g gin.IRouter) {
	g.POST("/addKey", svc.addKey)
	g.POST("/mapTo", svc.mapTo)
	g.GET("/keys", svc.getKeys)
	g.POST("/newApi", svc.newApi)
	g.POST("/testTransformer", svc.testTransformer)
}

func fakeUID() uint {
	return 1
}

func (svc *KeyStoreSvc) addKey(c *gin.Context) {
	req, aborted := utils.GinReqJson[apimodel.KeyStoreAddKeyPostRequest](c)
	if aborted {
		return
	}

	uid := fakeUID()

	key := dbmodel.OriginKey{
		Owner:    uid,
		Key:      req.Key,
		Provider: dbmodel.KeyProvider(req.Provider),
		EndPoint: req.GetEndpoint(),
	}
	err := shared.GetDB().Save(&key).Error
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, apimodel.KeyStoreAddKeyPost200Response{Id: int32(key.ID)})
}

func (svc *KeyStoreSvc) getKeys(c *gin.Context) {
	uid := fakeUID()
	var keys []dbmodel.OriginKey
	err := shared.GetDB().Where("owner = ?", uid).Find(&keys).Error
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	var res []apimodel.KeyStoreKeysGet200ResponseInner
	for _, key := range keys {
		ep := &apimodel.NullableString{}
		if key.EndPoint != "" {
			ep.Set(&key.EndPoint)
		}
		res = append(res, apimodel.KeyStoreKeysGet200ResponseInner{
			Id:       int32((key.ID)),
			Key:      key.Key,
			Provider: string(key.Provider),
			Endpoint: *ep,
		})
	}
	c.JSON(200, res)
}

func (svc *KeyStoreSvc) mapTo(c *gin.Context) {
	req, aborted := utils.GinReqJson[apimodel.KeyStoreMapToPostRequest](c)
	if aborted {
		return
	}

	// TODO: check if key exists & belongs to uid
	var rels []*dbmodel.MapOriginRelation
	for _, originId := range req.Keys {
		rels = append(rels, &dbmodel.MapOriginRelation{
			OriginID: uint(originId),
			MapID:    uint(req.MapTo),
		})
	}
	err := shared.GetDB().Save(rels).Error
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true})
}

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
		Owner: fakeUID(),
		Key:   key,
	}
	err := shared.GetDB().Save(&udk).Error
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, apimodel.NewKeyStoreNewApiPost200Response(int32(udk.ID), key))

}
