package data_loader

import (
	"github.com/dof-wiki/godof/parser"
	"pvftools/backend/common"
	"pvftools/backend/common/consts"
	"pvftools/backend/common/ctx"
	"pvftools/backend/common/log"
	"pvftools/backend/common/utils"
	"pvftools/backend/dao"
	"pvftools/backend/internal/data_source"
	"pvftools/backend/model"
	"strings"
)

type JobDataLoader struct{}

func (l *JobDataLoader) Name() string {
	return "职业"
}

func (l *JobDataLoader) CheckUpdate() bool {
	return CheckUpdateByMd5(consts.LstPathCharacter)
}

func (l *JobDataLoader) Load() {
	log.LogInfo("load job start.")
	defer log.LogInfo("load job end.")
	ds := data_source.GetDataSource()
	c, err := ds.GetFileContent(consts.LstPathCharacter)
	if err != nil {
		log.LogError("get file content error %v", err)
		return
	}
	p := parser.NewLstParser(c)
	dao.SaveFileMd5(consts.LstPathCharacter, utils.Md5(c))
	err = common.DB.Unscoped().Where("1=1").Delete(&model.Job{}).Error
	if err != nil {
		log.LogError("delete job error %v", err)
		return
	}
	jobs := make([]*model.Job, 0)
	for code, path := range p.GetPathMap() {
		realPath := "character/" + strings.ToLower(path)
		jc, err := ds.GetFileContent(realPath)
		if err != nil {
			log.LogError("get file content error %v", err)
			continue
		}
		jp := parser.NewParser(jc)
		jobName, err := jp.GetString("job")
		if err != nil {
			log.LogError("get jobname err %v", err)
			continue
		}
		groupTypeName, err := jp.GetStrings("growtype name")
		if err != nil {
			log.LogError("get group type name err %v", err)
			continue
		}

		for i, gtn := range groupTypeName {
			jobs = append(jobs, &model.Job{
				Code:         code,
				Job:          jobName[1 : len(jobName)-1],
				GrowType:     i,
				GrowTypeName: gtn,
			})
		}
	}
	if len(jobs) > 0 {
		common.DB.Create(&jobs)
	}

	ctx.Emit(consts.EventJobUpdate)
}

func init() {
	register("job", new(JobDataLoader))
}
