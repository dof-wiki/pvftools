package dao

import (
	"pvftools/backend/common"
	"pvftools/backend/model"
)

func GetEquipmentsByCodes(codes []int) ([]*model.Equipment, error) {
	equs := make([]*model.Equipment, 0, len(codes))
	if len(codes) == 0 {
		return equs, nil
	}
	if err := common.DB.Where("code in ?", codes).Find(&equs).Error; err != nil {
		return nil, err
	}
	return equs, nil
}

func GetEquipmentsByPaths(paths []string) ([]*model.Equipment, error) {
	equs := make([]*model.Equipment, 0, len(paths))
	if len(paths) == 0 {
		return equs, nil
	}
	if err := common.DB.Where("path in ?", paths).Find(&equs).Error; err != nil {
		return nil, err
	}
	return equs, nil
}

func GetEquipmentsByName(name string) ([]*model.Equipment, error) {
	equs := make([]*model.Equipment, 0)
	if err := common.DB.Where("name like ?", "%"+name+"%").Find(&equs).Error; err != nil {
		return nil, err
	}
	return equs, nil
}
