package data_loader

import (
	"github.com/dof-wiki/godof/parser"
	"github.com/dof-wiki/godof/parser/tree_parser"
	"github.com/samber/lo"
	"pvftools/backend/common"
	"pvftools/backend/common/consts"
	"pvftools/backend/common/log"
	"pvftools/backend/common/setting"
	"pvftools/backend/common/utils"
	"pvftools/backend/dao"
	"pvftools/backend/internal/data_source"
	"pvftools/backend/internal/progress"
	"pvftools/backend/model"
	"strings"
)

type NpcDataLoader struct {
	ch        chan *model.Npc
	stopCh    chan struct{}
	isRunning bool
	pc        *progress.ProgressController
}

func (l *NpcDataLoader) Name() string {
	return "NPC"
}

func (l *NpcDataLoader) CheckUpdate() bool {
	return CheckUpdateByMd5(consts.LstPathNpc)
}

func (l *NpcDataLoader) Load() {
	if l.isRunning {
		log.LogWarn("npc loader is running")
		return
	}
	l.isRunning = true
	l.ch = make(chan *model.Npc)
	l.stopCh = make(chan struct{})
	log.LogInfo("npc loader start")
	defer log.LogInfo("npc loader end")
	go l.run()

	err := common.DB.Unscoped().Where("1=1").Delete(new(model.Npc)).Error
	if err != nil {
		log.LogError("delete npc error %v", err)
		l.stopCh <- struct{}{}
		return
	}

	ds := data_source.GetDataSource()
	c, err := ds.GetFileContent(consts.LstPathNpc)
	if err != nil {
		log.LogError("get %s content err %v", consts.LstPathNpc, err)
		return
	}
	dao.SaveFileMd5(consts.LstPathNpc, utils.Md5(c))
	lst := parser.NewLstParser(c)
	codes := lo.Keys(lst.GetPathMap())
	l.pc = progress.NewProgressController("NPC", len(codes))
	utils.ConcurrentForEach(codes, func(code int) {
		defer func() {
			if err := recover(); err != nil {
				log.LogError("npc loader run %d-%s error: %v", code, lst.GetPath(code), err)
			}
		}()

		path := lst.GetPath(code)
		realPath := "npc/" + strings.ToLower(path)
		sc, err := ds.GetFileContent(realPath)
		if err != nil {
			log.LogError("get %s content err %v", path, err)
			return
		}
		p := tree_parser.NewTreeParser(sc)
		name := p.GetRoot().GetFirstChild(consts.LabelName).Value.GetString()
		npc := &model.Npc{
			ID:   code,
			Name: name,
			Path: realPath,
		}
		l.ch <- npc
	}, setting.UserSettings.ConcurrentCount)

	l.stopCh <- struct{}{}
}

func (l *NpcDataLoader) run() {
	defer func() {
		l.isRunning = false
	}()
	list := make([]*model.Npc, 0)
	for {
		select {
		case npc := <-l.ch:
			list = append(list, npc)
			l.pc.Increase(1)
		case <-l.stopCh:
			log.LogInfo("npc total: %d", len(list))
			if len(list) > 0 {
				if err := common.DB.Create(list).Error; err != nil {
					log.LogError("save npc err %v", err)
				}
			}
			l.pc.End()
			return
		}
	}
}

func init() {
	register("npc", new(NpcDataLoader))
}
