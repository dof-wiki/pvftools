package api

import (
	"pvftools/backend/common"
	"pvftools/backend/common/log"
	"pvftools/backend/model"
)

const (
	SearchTypeStackable = 1
	SearchTypeEqu       = 2
	SearchTypeSkill     = 3
	SearchTypeNpc       = 4
	SearchTypeDgn       = 5
)

type SearchResult struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ISearchModelList interface {
	ToSearchResult() []*SearchResult
}

type StackableList []*model.Stackable

func (l StackableList) ToSearchResult() []*SearchResult {
	list := make([]*SearchResult, 0)
	for _, v := range l {
		list = append(list, &SearchResult{
			Id:   v.Code,
			Name: v.Name,
		})
	}
	return list
}

type EquipmentList []*model.Equipment

func (l EquipmentList) ToSearchResult() []*SearchResult {
	list := make([]*SearchResult, 0)
	for _, v := range l {
		list = append(list, &SearchResult{
			Id:   v.Code,
			Name: v.Name,
		})
	}
	return list
}

type SkillList []*model.Skill

func (l SkillList) ToSearchResult() []*SearchResult {
	list := make([]*SearchResult, 0)
	for _, v := range l {
		list = append(list, &SearchResult{
			Id:   v.Code,
			Name: v.Name,
		})
	}
	return list
}

type NpcList []*model.Npc

func (l NpcList) ToSearchResult() []*SearchResult {
	list := make([]*SearchResult, 0)
	for _, v := range l {
		list = append(list, &SearchResult{
			Id:   v.ID,
			Name: v.Name,
		})
	}
	return list
}

type DungeonList []*model.Dungeon

func (l DungeonList) ToSearchResult() []*SearchResult {
	list := make([]*SearchResult, 0)
	for _, v := range l {
		list = append(list, &SearchResult{
			Id:   v.ID,
			Name: v.Name,
		})
	}
	return list
}

func (a *App) SearchByName(q string, tp int) []*SearchResult {
	var list any
	switch tp {
	case SearchTypeStackable:
		list = make(StackableList, 0)
	case SearchTypeEqu:
		list = make(EquipmentList, 0)
	case SearchTypeSkill:
		list = make(SkillList, 0)
	case SearchTypeNpc:
		list = make(NpcList, 0)
	case SearchTypeDgn:
		list = make(DungeonList, 0)
	}
	var err error
	if q == "" {
		err = common.DB.Find(&list).Error
	} else {
		err = common.DB.Where("name like ?", "%"+q+"%").Find(&list).Error
	}
	if err != nil {
		log.LogError("Search err %v", err)
		return nil
	}
	if list2, ok := list.(ISearchModelList); ok {
		return list2.ToSearchResult()
	} else {
		log.LogError("Search err")
	}
	return make([]*SearchResult, 0)
}
