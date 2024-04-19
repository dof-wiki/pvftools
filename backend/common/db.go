package common

import (
	"pvftools/backend/common/setting"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	var err error
	DB, err = gorm.Open(sqlite.Open(setting.DBFile()), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}
