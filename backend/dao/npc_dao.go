package dao

import (
	"pvftools/backend/common"
	"pvftools/backend/model"
)

func AllNpc() ([]*model.Npc, error) {
	list := make([]*model.Npc, 0)
	if err := common.DB.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}
