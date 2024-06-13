package model

import "pvftools/backend/common"

var models = []any{
	new(FileHash),
	new(Equipment),
	new(Job),
	new(Skill),
	new(Stackable),
	new(CustomAttrTmpl),
}

func init() {
	common.DB.AutoMigrate(models...)
}
