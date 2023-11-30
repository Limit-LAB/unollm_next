package dbmodel

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique,index"`
	Nickname string
	Hashed   string
}

type UserToken struct {
	Token     string `gorm:"varchar(128),primaryKey"`
	UserId    uint
	CreatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
