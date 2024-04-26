package box

import (
	"github.com/dof-wiki/godof/parser"
	"pvftools/backend/common/consts"
)

func GenerateRandomBox(params *BoxParams) *parser.Parser {
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

	subParser := parser.NewParser("")
	for _, group := range params.Items {
		vals := make([]int, 0)
		for _, item := range group.Items {
			vals = append(vals, item.Id, item.Rate, item.Count)
		}
		subParser.AddAny(consts.LabelEtc, parser.GenTokens(vals), true)
	}
	tokens := parser.GenTokenList(1)
	tokens = append(tokens, subParser.GetTokens()...)
	p.SetAny(consts.LabelBoosterInfo, tokens, true)
	return p
}
