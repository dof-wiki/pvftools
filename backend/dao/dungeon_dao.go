package dao

import (
	"pvftools/backend/common"
	"pvftools/backend/model"
)

func AllDungeon() ([]*model.Dungeon, error) {
	list := make([]*model.Dungeon, 0)
	if err := common.DB.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
