package dao

import (
	"pvftools/backend/common"
	"pvftools/backend/model"
)

func GetStackablesByCodes(codes []int) ([]*model.Stackable, error) {
	ret := make([]*model.Stackable, 0)
	if err := common.DB.Where("code in ?", codes).Find(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func GetStackablesByName(name string) ([]*model.Stackable, error) {
	ret := make([]*model.Stackable, 0)
	if err := common.DB.Where("name like ?", "%"+name+"%").Find(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}
