package dao

import (
	"pvftools/backend/common"
	"pvftools/backend/model"
)

func GetJobSkills(jobId int) ([]*model.Skill, error) {
	skills := make([]*model.Skill, 0)
	var err error
	if jobId >= 0 {
		err = common.DB.Where("job = ?", jobId).Find(&skills).Error
	} else {
		err = common.DB.Find(&skills).Error
	}
	return skills, err
}

func GetSkillByID(id int) (*model.Skill, error) {
	skill := &model.Skill{}
	err := common.DB.Where("id = ?", id).First(skill).Error
	return skill, err
}

func GetSkillByPath(path string) (*model.Skill, error) {
	skill := &model.Skill{}
	err := common.DB.Where("path = ?", path).First(skill).Error
	return skill, err
}
