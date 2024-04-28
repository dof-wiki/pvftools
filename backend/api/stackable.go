package api

import (
	"pvftools/backend/common/log"
	"pvftools/backend/dao"
)

func (a *App) GetItemsName(items []int) (map[int]string, error) {
	list, err := dao.GetStackablesByCodes(items)
	if err != nil {
		return nil, err
	}
	ret := make(map[int]string)
	for _, item := range list {
		ret[item.Code] = item.Name
	}
	return ret, nil
}

type SearchResult struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (a *App) SearchByName(q string) []*SearchResult {
	result := make([]*SearchResult, 0)
	items, err := dao.GetStackablesByName(q)
	if err != nil {
		log.LogError("Search stackable err %v", err)
	}
	for _, item := range items {
		result = append(result, &SearchResult{
			Id:   item.Code,
			Name: item.Name,
		})
	}
	equs, err := dao.GetEquipmentsByName(q)
	if err != nil {
		log.LogError("Search equipment err %v", err)
	}
	for _, equ := range equs {
		result = append(result, &SearchResult{
			Id:   equ.Code,
			Name: equ.Name,
		})
	}
	return result
}
