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

type DgnDataLoader struct {
	ch        chan *model.Dungeon
	stopCh    chan struct{}
	isRunning bool
	pc        *progress.ProgressController
}

func (l *DgnDataLoader) Name() string {
	return "副本"
}

func (l *DgnDataLoader) CheckUpdate() bool {
	return CheckUpdateByMd5(consts.LstPathDgn)
}

func (l *DgnDataLoader) Load() {
	if l.isRunning {
		log.LogWarn("dgn loader is running")
		return
	}
	l.isRunning = true
	l.ch = make(chan *model.Dungeon)
	l.stopCh = make(chan struct{})
	log.LogInfo("dgn loader start")
	defer log.LogInfo("dgn loader end")
	go l.run()

	err := common.DB.Unscoped().Where("1=1").Delete(new(model.Dungeon)).Error
	if err != nil {
		log.LogError("delete dungeon error %v", err)
		l.stopCh <- struct{}{}
		return
	}

	ds := data_source.GetDataSource()
	c, err := ds.GetFileContent(consts.LstPathDgn)
	if err != nil {
		log.LogError("get %s content err %v", consts.LstPathDgn, err)
		return
	}
	dao.SaveFileMd5(consts.LstPathDgn, utils.Md5(c))
	lst := parser.NewLstParser(c)
	codes := lo.Keys(lst.GetPathMap())
	l.pc = progress.NewProgressController("副本", len(codes))
	utils.ConcurrentForEach(codes, func(code int) {
		defer func() {
			if err := recover(); err != nil {
				log.LogError("dgn loader run error: %v", err)
			}
		}()

		path := lst.GetPath(code)
		realPath := "dungeon/" + strings.ToLower(path)
		sc, err := ds.GetFileContent(realPath)
		if err != nil {
			log.LogError("get %s content err %v", path, err)
			return
		}
		p := tree_parser.NewTreeParser(sc)
		name := p.GetRoot().GetFirstChild(consts.LabelName).Value.GetString()
		dgn := &model.Dungeon{
			ID:   code,
			Name: name,
			Path: realPath,
		}
		l.ch <- dgn
	}, setting.UserSettings.ConcurrentCount)

	l.stopCh <- struct{}{}
}

func (l *DgnDataLoader) run() {
	defer func() {
		l.isRunning = false
	}()
	list := make([]*model.Dungeon, 0)
	for {
		select {
		case dgn := <-l.ch:
			list = append(list, dgn)
			l.pc.Increase(1)
		case <-l.stopCh:
			log.LogInfo("dgn total: %d", len(list))
			if len(list) > 0 {
				if err := common.DB.Create(list).Error; err != nil {
					log.LogError("save dgn err %v", err)
				}
			}
			l.pc.End()
			return
		}
	}
}

func init() {
	register("dgn", new(DgnDataLoader))
}
