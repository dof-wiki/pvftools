package dao

import (
	"pvftools/backend/common"
	"pvftools/backend/common/consts"
	"pvftools/backend/model"
)

func AllJobs() ([]*model.Job, error) {
	jobs := make([]*model.Job, 0)
	if err := common.DB.Find(&jobs).Error; err != nil {
		return nil, err
	}
	for _, job := range jobs {
		job.JobName = consts.JobChineseName[job.Job]
	}
	return jobs, nil
}
