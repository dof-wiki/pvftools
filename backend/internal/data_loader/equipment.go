package data_loader

import (
	"github.com/dof-wiki/godof/parser"
	"github.com/samber/lo"
	"pvftools/backend/common"
	"pvftools/backend/common/consts"
	"pvftools/backend/common/log"
	"pvftools/backend/common/setting"
	"pvftools/backend/common/utils"
	"pvftools/backend/dao"
	"pvftools/backend/internal/data_source"
	"pvftools/backend/model"
)

type EquipmentDataLoader struct {
	equCh     chan *model.Equipment
	stopCh    chan struct{}
	isRunning bool
}

func (l *EquipmentDataLoader) Name() string {
	return "装备"
}

func (l *EquipmentDataLoader) CheckUpdate() bool {
	return CheckUpdateByMd5(consts.LstPathEquipment)
}

func (l *EquipmentDataLoader) Load() {
	if l.isRunning {
		log.LogWarn("equipment loader is running")
		return
	}
	l.isRunning = true
	log.LogInfo("load equipment start.")
	defer log.LogInfo("load equipment end.")
	go l.run()

	ds := data_source.GetDataSource()
	c, err := ds.GetFileContent(consts.LstPathEquipment)
	if err != nil {
		log.LogError("get %s content err %v", consts.LstPathEquipment, err)
		return
	}

	p := parser.NewLstParser(c)
	dao.SaveFileMd5(consts.LstPathEquipment, utils.Md5(c))
	err = common.DB.Unscoped().Where("1=1").Delete(&model.Equipment{}).Error
	if err != nil {
		log.LogError("delete equ error %v", err)
		return
	}

	codes := lo.Keys(p.GetPathMap())
	utils.ConcurrentForEach(codes, func(code int) {
		defer func() {
			if r := recover(); r != nil {
				log.LogError("error is %v", r)
			}
		}()
		path := p.GetPath(code)
		realPath := "equipment/" + path
		content, err := ds.GetFileContent(realPath)
		if err != nil {
			log.LogError("get %s content err %v", realPath, err)
			return
		}
		fp := parser.NewParser(content)
		name, err := fp.GetString("name")
		if err != nil {
			log.LogError("get equipment name err %v", err)
			return
		}
		tv := fp.GetAny("equipment type")
		equTp := utils.RemoveWrapper(tv.GetCleanTokens()[0].RawContent())
		rarity, _ := fp.GetInt("rarity")
		partset, _ := fp.GetInt("part set index")
		miniLevel, _ := fp.GetInt("minimum level")
		equ := &model.Equipment{
			Code:      code,
			Name:      name,
			Path:      realPath,
			Rarity:    rarity,
			Type:      equTp,
			PartSet:   partset,
			MiniLevel: miniLevel,
		}
		l.equCh <- equ
	}, setting.UserSettings.ConcurrentCount)

	l.stopCh <- struct{}{}
}

func (l *EquipmentDataLoader) run() {
	defer func() {
		l.isRunning = false
	}()
	equList := make([]*model.Equipment, 0, 1000)
	for {
		select {
		case <-l.stopCh:
			if len(equList) > 0 {
				err := common.DB.Create(equList).Error
				if err != nil {
					log.LogError("save equipment err %v", err)
					return
				}
			}
			return
		case equ := <-l.equCh:
			equList = append(equList, equ)
			if len(equList) >= 1000 {
				err := common.DB.Create(equList).Error
				if err != nil {
					log.LogError("save equipment err %v", err)
					return
				}
				equList = make([]*model.Equipment, 0, 1000)
			}
		}
	}
}

func init() {
	register("equipment", new(EquipmentDataLoader))
}
