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

type StackableDataLoader struct {
	stackableCh chan *model.Stackable
	stopCh      chan struct{}
	isRunning   bool
}

func (l *StackableDataLoader) Name() string {
	return "道具"
}

func (l *StackableDataLoader) CheckUpdate() bool {
	return CheckUpdateByMd5(consts.LstPathStackable)
}

func (l *StackableDataLoader) Load() {
	if l.isRunning {
		log.LogWarn("stackable loader is running")
		return
	}
	l.isRunning = true
	log.LogInfo("stackable loader start")
	defer log.LogInfo("stackable loader end")
	go l.run()

	ds := data_source.GetDataSource()
	c, err := ds.GetFileContent(consts.LstPathStackable)
	if err != nil {
		log.LogError("get stackable lst content error: %v", err)
		return
	}
	lst := parser.NewLstParser(c)
	dao.SaveFileMd5(consts.LstPathStackable, utils.Md5(c))
	err = common.DB.Unscoped().Where("1=1").Delete(new(model.Stackable)).Error
	if err != nil {
		log.LogError("delete stackable err %v", err)
		return
	}
	codes := lo.Keys(lst.GetPathMap())
	utils.ConcurrentForEach(codes, func(code int) {
		defer func() {
			if err := recover(); err != nil {
				log.LogError("stackable loader run error: %v", err)
			}
		}()
		path := lst.GetPath(code)
		realPath := "stackable/" + path
		content, err := ds.GetFileContent(realPath)
		if err != nil {
			log.LogError("get stackable content error: %v", err)
			return
		}
		p := parser.NewParser(content)
		name, _ := p.GetString("name")
		stackable := &model.Stackable{
			Code: code,
			Name: name,
			Path: realPath,
		}
		l.stackableCh <- stackable
	}, setting.UserSettings.ConcurrentCount)

	l.stopCh <- struct{}{}
}

func (l *StackableDataLoader) run() {
	defer func() {
		l.isRunning = false
	}()
	list := make([]*model.Stackable, 0, 1000)
	for {
		select {
		case <-l.stopCh:
			if len(list) > 0 {
				err := common.DB.Create(list).Error
				if err != nil {
					log.LogError("save stackable err %v", err)
				}
			}
			return
		case stk := <-l.stackableCh:
			list = append(list, stk)
			if len(list) >= 1000 {
				err := common.DB.Create(list).Error
				if err != nil {
					log.LogError("save stackable err %v", err)
				}
				list = make([]*model.Stackable, 0, 1000)
			}
		}
	}
}

func init() {
	register("stackable", new(StackableDataLoader))
}
