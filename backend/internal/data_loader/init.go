package data_loader

import (
	"pvftools/backend/common/log"
	"pvftools/backend/common/utils"
	"pvftools/backend/dao"
	"pvftools/backend/internal/data_source"
)

type ILoader interface {
	Name() string
	CheckUpdate() bool
	Load()
}

func CheckUpdateByMd5(path string) bool {
	cur := dao.GetFileMd5(path)
	if cur == "" {
		return true
	}
	c, err := data_source.GetDataSource().GetFileContent(path)
	if err != nil {
		log.LogError("get file content error: %v", err)
		return false
	}
	return utils.Md5(c) != cur
}
