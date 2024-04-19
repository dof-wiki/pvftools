package api

import "pvftools/backend/internal/data_source"

func (a *App) GetFileContent(path string) (string, error) {
	return data_source.GetDataSource().GetFileContent(path)
}

func (a *App) SaveFile(path, content string) error {
	return data_source.GetDataSource().SaveFileContent(path, content)
}
