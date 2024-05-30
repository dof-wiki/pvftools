package proto

import (
	"github.com/dof-wiki/godof/parser/tree_parser"
	"pvftools/backend/common/consts"
)

type BreathCheckSkillItem struct {
	Index         int               `json:"index"`
	IncrementList []*SkillIncrement `json:"increment_list"`
	Explain       string            `json:"explain"`
}

func (s *BreathCheckSkillItem) SkillValues() tree_parser.Value {
	vals := []any{s.Index}
	for _, inc := range s.IncrementList {
		vals = append(vals, "\n", "\t", inc.Job, inc.SkillId, inc.Type, inc.DataType, inc.DataIndex, inc.IncrementType, inc.Value)
	}
	return tree_parser.GenTokenList(vals...)
}

type BreathCheckItem struct {
	Job       int                     `json:"job"`
	SubJob    int                     `json:"sub_job"`
	EquType   string                  `json:"equ_type"`
	SkillList []*BreathCheckSkillItem `json:"skill_list"`
}

func (c *BreathCheckItem) TreeNode() *tree_parser.TreeNode {
	checkSubNode := tree_parser.NewTreeNode("root", true)
	for _, skl := range c.SkillList {
		checkSubNode.AddChild(tree_parser.NewTreeNode(consts.LabelSkill, true, skl.SkillValues()...))
		checkSubNode.AddChild(tree_parser.NewTreeNode(consts.LabelSkillExplain, false, tree_parser.GenTokens(skl.Explain)...))
	}
	sub := checkSubNode.Render()
	values := tree_parser.GenTokenList(c.Job, c.SubJob, c.EquType)
	values = append(values, tree_parser.NewDelimiterToken("\n"), tree_parser.NewRawContentToken(sub))
	n := tree_parser.NewTreeNode(consts.LabelCheck, true, values...)
	return n
}

type BreathData struct {
	Id           int                `json:"id"`
	Name         string             `json:"name"`
	Probability1 int                `json:"probability1"`
	Probability2 int                `json:"probability2"`
	CheckList    []*BreathCheckItem `json:"check_list"`
	HasChange    bool               `json:"has_change"`
}
