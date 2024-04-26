package box

import (
	"github.com/dof-wiki/godof/parser"
	"pvftools/backend/common/consts"
)

func GenerateSelectionBox(params *BoxParams) *parser.Parser {
	p := parser.NewParser("")
	p.SetString(consts.LabelName, params.Name)
	p.SetString(consts.LabelExplain, params.Explain)
	p.SetInt(consts.LabelGrade, params.Grade)
	p.SetInt(consts.LabelRarity, params.Rarity)
	p.SetString(consts.LabelUsableJob, "[all]", true)
	p.SetString(consts.LabelAttachType, "[free]")
	p.SetInt(consts.LabelMinimumLevel, 1)
	p.SetAny(consts.LabelIcon, parser.GenTokenList("item/stackable/consumption_cn.img", 179))
	p.SetAny(consts.LabelIconMark, parser.GenTokenList("item/IconMark.img", 64))
	p.SetAny(consts.LabelStackableType, parser.GenTokenList("[booster]", 0))
	p.SetString(consts.LabelMoveWav, "BONE_TOUCH")
	p.SetInt(consts.LabelStackLimit, 1000)
	if len(params.Items) == 1 {
		p.SetInt(consts.LabelBoosterCategoryNum, 0)
	} else {
		p.SetInts(consts.LabelBoosterCategoryNum, []int{1, len(params.Items)})
	}
	p.SetInt(consts.LabelBoosterSelectionNum, 1)

	for i, group := range params.Items {
		subParser := parser.NewParser("")
		vals := make([]int, 0)
		for _, item := range group.Items {
			vals = append(vals, item.Id, item.Count)
		}
		subParser.AddAny(params.ItemType, parser.GenTokens(vals), true)
		p.AddAny(consts.LabelBoosterSelectCategory, append(parser.GenTokenList(i, 0), subParser.GetTokens()...), true)
	}

	if len(params.CategoryNames) > 0 {
		p.SetStrings(consts.LabelBoosterCategoryName, params.CategoryNames, true)
	}
	return p
}
