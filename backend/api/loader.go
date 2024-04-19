package api

import "pvftools/backend/internal/data_loader"

func (a *App) CheckDataUpdate() []*data_loader.CheckUpdateResult {
	return data_loader.DataLoaderMgr.CheckUpdate()
}

func (a *App) LoadData(keys []string) {
	go data_loader.DataLoaderMgr.LoadData(keys...)
}
