package data_loader

import (
	"github.com/dof-wiki/godof/parser"
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

type SkillDataLoader struct {
	skillCh   chan *model.Skill
	stopCh    chan struct{}
	isRunning bool
	pc        *progress.ProgressController
}

func (l *SkillDataLoader) Name() string {
	return "技能"
}

func (l *SkillDataLoader) CheckUpdate() bool {
	ds := data_source.GetDataSource()
	c, err := ds.GetFileContent(consts.LstPathSkill)
	if err != nil {
		log.LogError("get %s content err %v", consts.LstPathSkill, err)
		return false
	}
	lst := parser.NewLstParser(c)
	for _, path := range lst.GetPathMap() {
		if CheckUpdateByMd5("skill/" + strings.ToLower(path)) {
			return true
		}
	}
	return false
}

func (l *SkillDataLoader) Load() {
	if l.isRunning {
		log.LogWarn("skill loader is running")
		return
	}
	l.isRunning = true
	l.skillCh = make(chan *model.Skill)
	l.stopCh = make(chan struct{})
	log.LogInfo("skill loader start")
	defer log.LogInfo("skill loader end")
	go l.run()
	ds := data_source.GetDataSource()
	c, err := ds.GetFileContent(consts.LstPathSkill)
	if err != nil {
		log.LogError("get %s content err %v", consts.LstPathSkill, err)
		l.stopCh <- struct{}{}
		return
	}
	lst := parser.NewLstParser(c)
	type SkillItem struct {
		jobId int
		code  int
		path  string
	}
	items := make([]*SkillItem, 0)
	for jobId, path := range lst.GetPathMap() {
		realPath := "skill/" + strings.ToLower(path)
		sc, err := ds.GetFileContent(realPath)
		if err != nil {
			log.LogError("get %s content err %v", path, err)
			continue
		}
		subLst := parser.NewLstParser(sc)
		dao.SaveFileMd5(realPath, utils.Md5(sc))
		for code, v := range subLst.GetPathMap() {
			items = append(items, &SkillItem{
				jobId: jobId,
				code:  code,
				path:  "skill/" + strings.ToLower(v),
			})
		}
	}
	err = common.DB.Unscoped().Where("1=1").Delete(&model.Skill{}).Error
	if err != nil {
		log.LogError("delete skill error %v", err)
		l.stopCh <- struct{}{}
		return
	}
	l.pc = progress.NewProgressController("技能", len(items))
	utils.ConcurrentForEach(items, func(item *SkillItem) {
		c, err := ds.GetFileContent(item.path)
		if err != nil {
			log.LogError("get %s content err %v", item.path, err)
			return
		}
		p := parser.NewParser(c)
		name, _ := p.GetString("name")
		skill := &model.Skill{
			Job:  item.jobId,
			Code: item.code,
			Name: name,
			Path: item.path,
		}
		l.skillCh <- skill
	}, setting.UserSettings.ConcurrentCount)
	l.stopCh <- struct{}{}
}

func (l *SkillDataLoader) run() {
	defer func() {
		l.isRunning = false
	}()
	skillList := make([]*model.Skill, 0)
	for {
		select {
		case skill := <-l.skillCh:
			skillList = append(skillList, skill)
			l.pc.Increase(1)
		case <-l.stopCh:
			log.LogInfo("skill total: %d", len(skillList))
			if len(skillList) > 0 {
				if err := common.DB.Create(skillList).Error; err != nil {
					log.LogError("save skill err %v", err)
				}
			}
			l.pc.End()
			return
		}
	}
}

func init() {
	register("skill", new(SkillDataLoader))
}
