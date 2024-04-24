package api

import (
	"pvftools/backend/dao"
	"pvftools/backend/model"
)

func (a *App) GetJobInfo() ([]*model.Job, error) {
	return dao.AllJobs()
}
