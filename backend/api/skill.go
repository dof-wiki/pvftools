package api

import (
	"fmt"
	"github.com/dof-wiki/godof/parser"
	"github.com/samber/lo"
	"pvftools/backend/common/consts"
	"pvftools/backend/dao"
	"pvftools/backend/model"
	"strings"
)

func (a *App) GetJobSkill(job int) ([]*model.Skill, error) {
	return dao.GetJobSkills(job)
}

type SkillDescription struct {
	Tp   int     `json:"tp"`
	Idx  int     `json:"idx"`
	Rate float64 `json:"rate"`
}

type SkillData struct {
	Name                string              `json:"name"`
	Code                int                 `json:"code"`
	JobStr              string              `json:"job_str"`
	StaticData          []int               `json:"static_data"`
	LevelData           [][]int             `json:"level_data"`
	LevelDataNum        int                 `json:"level_data_num"`
	Description         []*SkillDescription `json:"description"`
	DescriptionTemplate string              `json:"description_template"`
}

func (a *App) GetSkillDetail(id int) (*SkillData, error) {
	m, err := dao.GetSkillByID(id)
	if err != nil {
		return nil, err
	}
	p := a.getParser(m.Path)
	sklData := new(SkillData)
	sklData.Code = m.Code
	sklData.JobStr = consts.JobID2Name[m.Job]

	dungeon := p.GetAny(consts.LabelDungeon)
	var levelInfo, staticData []int
	if dungeon == nil {
		levelInfo, _ = p.GetInts(consts.LabelLevelInfo)
		staticData, _ = p.GetInts(consts.LabelStaticData)
	} else {
		levelInfo, _ = dungeon.GetSub(consts.LabelLevelInfo).GetInts()
		staticData, _ = dungeon.GetSub(consts.LabelStaticData).GetInts()
	}
	sklData.StaticData = staticData
	if len(levelInfo) > 0 {
		sklData.LevelDataNum = levelInfo[0]
		sklData.LevelData = lo.Chunk(levelInfo[1:], sklData.LevelDataNum)
	}

	descData := p.GetAny(consts.LabelLevelProperty)
	if descData != nil && len(descData.GetCleanTokens()) > 3 {
		numberTokens := make([]*parser.Token, 0)
		desc := ""
		for _, v := range descData.GetCleanTokens() {
			if v.IsInt() || v.IsFloat() {
				numberTokens = append(numberTokens, v)
			} else if v.IsString() {
				desc = v.RawContent()
			}
		}
		sklData.Description = make([]*SkillDescription, 0)
		for _, items := range lo.Chunk(numberTokens[2:], 3) {
			sklData.Description = append(sklData.Description, &SkillDescription{
				Tp:   items[0].GetInt(),
				Idx:  items[1].GetInt(),
				Rate: items[2].GetFloat(),
			})
		}
		sklData.DescriptionTemplate = desc
	}
	return sklData, nil
}

func (a *App) setSkillData(path string, static []int, level [][]int) error {
	p := a.getParser(path)
	dungeon := p.GetAny(consts.LabelDungeon)
	if dungeon == nil {
		p.SetInts(consts.LabelStaticData, static, true)
		curLevelInfo, _ := p.GetInts(consts.LabelLevelInfo)
		if len(curLevelInfo) > 0 {
			num := curLevelInfo[0]
			p.SetInts(consts.LabelLevelInfo, append([]int{num}, lo.Flatten(level)...), true)
		}
	} else {
		dungeon.SetSub(consts.LabelStaticData, parser.GenTokens(static), true)
		curLevelInfo, _ := dungeon.GetSub(consts.LabelLevelInfo).GetInts()
		if len(curLevelInfo) > 0 {
			num := curLevelInfo[0]
			dungeon.SetSub(consts.LabelLevelInfo, parser.GenTokens(append([]int{num}, lo.Flatten(level)...)), true)
		}
		p.SetAny(consts.LabelDungeon, dungeon.GetFull(), true)
	}
	return a.saveData(path)
}

func (a *App) SetSkillData(id int, static []int, level [][]int, withTp bool) error {
	skl, err := dao.GetSkillByID(id)
	if err != nil {
		return err
	}
	err = a.setSkillData(skl.Path, static, level)
	if err != nil {
		return err
	}
	if withTp {
		prefix := strings.TrimSuffix(skl.Path, ".skl")
		tpPath := fmt.Sprintf("%sex.skl", prefix)
		tp, err := dao.GetSkillByPath(tpPath)
		if err != nil {
			return err
		}
		if tp != nil {
			err = a.setSkillData(tp.Path, static, level)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
