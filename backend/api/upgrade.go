package api

import (
	"pvftools/backend/common/consts"
	"pvftools/backend/common/log"
	"pvftools/backend/internal/data_source"
	"pvftools/backend/internal/upgrade"
	"pvftools/backend/proto"
	"sync"
)

var upgradeController = make(map[string]*upgrade.UpgradeController)
var mu sync.Mutex

func getUpgradeController(path string) *upgrade.UpgradeController {
	mu.Lock()
	defer mu.Unlock()
	if c, ok := upgradeController[path]; ok {
		return c
	}
	c := upgrade.NewUpgradeController(path)
	upgradeController[path] = c
	return c
}

func (a *App) GetUpgradeData(tp int) []*proto.UpgradeItem {
	var path string
	if tp == 1 {
		path = consts.PathUpgrade
	} else {
		path = consts.PathAmplifyUpgrade
	}
	c := getUpgradeController(path)
	return c.GetItems()
}

func (a *App) SetUpgradeItem(tp int, items []*proto.UpgradeItem) {
	var path string
	if tp == 1 {
		path = consts.PathUpgrade
	} else {
		path = consts.PathAmplifyUpgrade
	}
	c := getUpgradeController(path)
	c.UpdateItems(items)
}

func (a *App) SaveUpgradeData(tp int) {
	var path string
	if tp == 1 {
		path = consts.PathUpgrade
	} else {
		path = consts.PathAmplifyUpgrade
	}
	c := getUpgradeController(path)
	content := c.Render()
	if err := data_source.GetDataSource().SaveFileContent(path, content); err != nil {
		log.LogError("保存失败 %v", err)
	}
}
