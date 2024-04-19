package data_source

import (
	"pvftools/backend/common/setting"

	"github.com/dof-wiki/godof/pvf_source"
)

type IDataSource interface {
	GetFileContent(path string) (string, error)
	SaveFileContent(path, content string) error
}

var dataSource IDataSource

func createDataSource() IDataSource {
	if setting.UserSettings.Mode == setting.ModePvfUtility {
		return pvf_source.NewPvfUtilitySource(setting.UserSettings.Target)
	} else {
		return pvf_source.NewFileSystemSource(setting.UserSettings.Target)
	}
}

func GetDataSource() IDataSource {
	if dataSource == nil {
		dataSource = createDataSource()
	}
	return dataSource
}
