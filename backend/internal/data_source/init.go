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
var dataMode string

func createDataSource() IDataSource {
	dataMode = setting.UserSettings.Mode
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
	if dataMode != setting.UserSettings.Mode {
		dataSource = createDataSource()
	}
	return dataSource
}
