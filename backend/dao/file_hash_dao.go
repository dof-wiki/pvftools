package dao

import (
	"pvftools/backend/common"
	"pvftools/backend/model"
)

func GetFileMd5(path string) string {
	m := new(model.FileHash)
	if err := common.DB.Where("path = ?", path).First(m).Error; err != nil {
		return ""
	}
	return m.Md5
}

func SaveFileMd5(path, md5 string) {
	m := new(model.FileHash)
	m.Path = path
	common.DB.Where("path = ?", path).First(m)
	m.Md5 = md5
	common.DB.Save(m)
}
