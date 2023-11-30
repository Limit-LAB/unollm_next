package apimodel

import "gorm.io/gorm"

type OriginKey struct {
	gorm.Model
	Key      string      `gorm:"text"`
	Owner    uint        `gorm:"index"`
	Provider KeyProvider `gorm:"enum:chatglm,chatgpt"`
	EndPoint string      `gorm:"text"`
	Disabled bool        `gorm:"default:false"`
}

type KeyProvider string

const (
	KeyProviderChatGLM KeyProvider = "chatglm"
	KeyProviderChatGPT KeyProvider = "chatgpt"
)

type KeyGroup struct {
	gorm.Model
	Owner uint `gorm:"index"`
}
