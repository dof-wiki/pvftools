package upgrade

import (
	"github.com/dof-wiki/godof/parser/tree_parser"
	"github.com/samber/lo"
	"pvftools/backend/common/consts"
	"pvftools/backend/internal/data_source"
	"pvftools/backend/proto"
)

/*
2>0.018>0.1>0.033>0.3>0.6>0>1>0>3037> 10> 1>0.018>0.1>0.033>0.17> 0.6

17一组
0: 等级
6: 失败概率, x / 10000
7: 失败处理, 1-不处理, 2-降级(参数决定降多少级), 3-破坏
8: 失败处理参数
9: 材料
10: 材料数量
*/

const (
	IndexFailRate    = 6
	IndexFailOp      = 7
	IndexFailOpValue = 8
	IndexCostId      = 9
	IndexCostCount   = 10
)

type UpgradeController struct {
	p *tree_parser.TreeParser

	items [][]any
}

func NewUpgradeController(path string) *UpgradeController {
	content, _ := data_source.GetDataSource().GetFileContent(path)
	c := &UpgradeController{
		p:     tree_parser.NewTreeParser(content),
		items: make([][]any, 0),
	}
	c.parse()
	return c
}

func (c *UpgradeController) parse() {
	root := c.p.GetRoot()
	table := root.GetFirstChild(consts.LabelTable)
	c.items = lo.Chunk(table.Value.Get(), 17)
}

func (c *UpgradeController) GetItems() []*proto.UpgradeItem {
	items := make([]*proto.UpgradeItem, 0)
	for idx, l := range c.items {
		level := idx + 1
		if level > 31 {
			break
		}
		items = append(items, &proto.UpgradeItem{
			Level:       level,
			FailRate:    l[IndexFailRate].(int),
			FailOp:      l[IndexFailOp].(int),
			FailOpValue: l[IndexFailOpValue].(int),
			CostId:      l[IndexCostId].(int),
			CostCount:   l[IndexCostCount].(int),
		})
	}
	return items
}

func (c *UpgradeController) Render() string {
	table := c.p.GetRoot().GetFirstChild(consts.LabelTable)
	args := append([]any{"\n", "\t"}, lo.Flatten(c.items)...)
	args = append(args, "\n")
	table.Value = tree_parser.GenTokenList(args...)
	return c.p.Render()
}

func (c *UpgradeController) UpdateItems(items []*proto.UpgradeItem) {
	for _, item := range items {
		idx := item.Level - 1
		c.items[idx][IndexFailRate] = item.FailRate
		c.items[idx][IndexFailOp] = item.FailOp
		c.items[idx][IndexFailOpValue] = item.FailOpValue
		c.items[idx][IndexCostId] = item.CostId
		c.items[idx][IndexCostCount] = item.CostCount
	}
}
