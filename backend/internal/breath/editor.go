package breath

import (
	"github.com/dof-wiki/godof/parser/tree_parser"
	"github.com/samber/lo"
	"pvftools/backend/common/consts"
	"pvftools/backend/proto"
)

type BreathEditor struct {
	p *tree_parser.TreeParser

	Data *proto.BreathData
}

func NewBreathEditor(id int, p *tree_parser.TreeParser) *BreathEditor {
	e := &BreathEditor{
		p: p,
		Data: &proto.BreathData{
			Id: id,
		},
	}
	e.parse()
	return e
}

func (b *BreathEditor) parse() {
	root := b.p.GetRoot()
	b.Data.Name = root.GetFirstChild(consts.LabelName).Value.GetString()
	enchant := root.GetFirstChild(consts.Label3ChoroEnchant)
	probabilities := enchant.GetFirstChild(consts.LabelProbability).Value.GetInts()
	b.Data.Probability1 = probabilities[0]
	b.Data.Probability2 = probabilities[1]

	checkNodes := enchant.GetChildren(consts.LabelCheck)
	checkItems := make([]*proto.BreathCheckItem, 0, len(checkNodes))
	for _, c := range checkNodes {
		idx := 0
		for i, v := range c.Value {
			if v.IsKey() {
				idx = i
				break
			}
		}
		values := c.Value[idx:]
		vals := c.Value[:20].Get()
		job := vals[0].(int)
		subJob := vals[1].(int)
		equType := vals[2].(string)
		checkNode := tree_parser.NewTreeNode(consts.LabelCheck, true, values...)
		skillNodes := checkNode.GetChildren(consts.LabelSkill)
		skillExplainNodes := checkNode.GetChildren(consts.LabelSkillExplain)
		skillList := make([]*proto.BreathCheckSkillItem, 0, len(skillNodes))
		for i, s := range skillNodes {
			skillExplain := skillExplainNodes[i].Value.GetString()
			data := s.Value.Get()
			skillIncrData := lo.Chunk(data[1:], 7)
			skillList = append(skillList, &proto.BreathCheckSkillItem{
				Index: data[0].(int),
				IncrementList: lo.Map(skillIncrData, func(d []any, _ int) *proto.SkillIncrement {
					item := new(proto.SkillIncrement)
					item.FromData([7]any(d))
					return item
				}),
				Explain: skillExplain,
			})
		}
		checkItems = append(checkItems, &proto.BreathCheckItem{
			Job:       job,
			SubJob:    subJob,
			EquType:   equType,
			SkillList: skillList,
		})
	}
	b.Data.CheckList = checkItems
}

func (b *BreathEditor) SetSkill(job, subJob int, equType string, skillList []*proto.BreathCheckSkillItem) {
	for _, v := range b.Data.CheckList {
		if v.Job == job && v.SubJob == subJob && v.EquType == equType {
			v.SkillList = skillList
			return
		}
	}
	b.Data.CheckList = append(b.Data.CheckList, &proto.BreathCheckItem{
		Job:       job,
		SubJob:    subJob,
		EquType:   equType,
		SkillList: skillList,
	})
	b.Data.HasChange = true
}

func (b *BreathEditor) Render() string {
	root := b.p.GetRoot()
	enchant := root.GetFirstChild(consts.Label3ChoroEnchant)
	probablityNode := enchant.GetFirstChild(consts.LabelProbability)
	probablityNode.Value = tree_parser.GenTokenList("\n", "\t", b.Data.Probability1, b.Data.Probability2, "\n", "\t")
	children := []*tree_parser.TreeNode{probablityNode}
	for _, c := range b.Data.CheckList {
		children = append(children, c.TreeNode())
	}
	enchant.SetChildren(children)
	return b.p.Render()
}
