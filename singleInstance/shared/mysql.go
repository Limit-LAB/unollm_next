package shared

import (
	"go.limit.dev/unollm/singleInstance/model/dbmodel"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _db *gorm.DB

func migrate() error {
	return _db.AutoMigrate(
		&dbmodel.MapOriginRelation{},
		&dbmodel.OriginKey{},
		&dbmodel.UserDefinedKey{},
		&dbmodel.User{},
	)
}

func InitMySql() error {
	var err error
	_db, err = gorm.Open(
		mysql.Open(GetCfg().SQL),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
			PrepareStmt:                              true,
		})
	if err != nil {
		return err
	}
	return migrate()
}

func GetDB() *gorm.DB {
	return _db
}
