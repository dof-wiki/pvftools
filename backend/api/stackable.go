package api

import "pvftools/backend/dao"

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
