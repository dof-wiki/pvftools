package api

import (
	"errors"
	"pvftools/backend/common/consts"
	"pvftools/backend/common/log"
	"pvftools/backend/internal/breath"
	"pvftools/backend/internal/data_source"
	"pvftools/backend/proto"
)

var breathPath = map[int]string{
	1254: consts.PathBreathRed,
	1255: consts.PathBreathBlue,
	1256: consts.PathBreathGreen,
	1257: consts.PathBreathYellow,
}

var breathEditor = make(map[int]*breath.BreathEditor)

func (a *App) getBreathEditor(id int) *breath.BreathEditor {
	cache := breathEditor[id]
	if cache != nil {
		return cache
	}
	p := a.getTreeParser(breathPath[id])
	editor := breath.NewBreathEditor(id, p)
	breathEditor[id] = editor
	return editor
}

func (a *App) GetBreathData(id int) (*proto.BreathData, error) {
	editor := a.getBreathEditor(id)
	return editor.Data, nil
}

func (a *App) SetBreathBaseData(id, probability1, probability2 int) {
	editor := a.getBreathEditor(id)
	editor.Data.Probability1 = probability1
	editor.Data.Probability2 = probability2
	editor.Data.HasChange = true
}

func (a *App) SetBreathSkills(id, job, subJob int, equType string, skillList []*proto.BreathCheckSkillItem) {
	editor := a.getBreathEditor(id)
	editor.SetSkill(job, subJob, equType, skillList)
}

func (a *App) SaveBreath(id int) error {
	editor := a.getBreathEditor(id)
	if p, ok := breathPath[id]; ok {
		//p = "2.stk"
		if err := data_source.GetDataSource().SaveFileContent(p, editor.Render()); err != nil {
			log.LogError("save data %s err %v", p, err)
			return err
		} else {
			editor.Data.HasChange = false
			return nil
		}
	} else {
		return errors.New("id error")
	}
}

func (a *App) ReloadBreath(id int) *proto.BreathData {
	a.delTreeParser(breathPath[id])
	p := a.getTreeParser(breathPath[id])
	editor := breath.NewBreathEditor(id, p)
	breathEditor[id] = editor
	return editor.Data
}
