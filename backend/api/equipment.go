package api

import (
	"fmt"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"os"
	"path"
	"pvftools/backend/common"
	"pvftools/backend/common/consts"
	"pvftools/backend/common/log"
	"pvftools/backend/internal/script"
	"pvftools/backend/model"
	"strings"
	"time"
)

type EquQueryCond struct {
	Rarity      []int    `json:"rarity"`
	Types       []string `json:"types"`
	PartSetType int      `json:"part_set_type"`
}

func (a *App) QueryEquipments(cond EquQueryCond) ([]*model.Equipment, error) {
	equs := make([]*model.Equipment, 0)
	session := common.DB.Where("rarity in ?", cond.Rarity).Where("type in ?", cond.Types)
	switch cond.PartSetType {
	case 1:
		session = session.Where("part_set > 0")
	case 2:
		session = session.Where("part_set = 0")
	}
	if err := session.Find(&equs).Error; err != nil {
		return nil, err
	}
	return equs, nil
}

func (a *App) EquipmentBatchHandle(idList []int, actionList []*script.Action) error {
	lst := a.getLstParser(consts.LstPathEquipment)
	for _, code := range idList {
		realPath := "equipment/" + lst.GetPath(code)
		p := a.getParser(realPath)
		for _, act := range actionList {
			act.Run(p)
		}
		if err := a.saveData(realPath); err != nil {
			log.LogError("%s保存失败", realPath)
		}
	}
	return nil
}

func (a *App) ExportFiles(paths []string) error {
	filename := fmt.Sprintf("export_%d.txt", time.Now().Unix())
	dirPath, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		CanCreateDirectories: true,
	})
	if err != nil {
		return err
	}

	filepath := path.Join(dirPath, filename)
	data := strings.Join(paths, "\r\n")
	fmt.Println(data)
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(data)
	return err
}
