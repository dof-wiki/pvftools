package world_drop

import (
	"github.com/dof-wiki/godof/parser"
	"github.com/samber/lo"
	"pvftools/backend/common/consts"
	"pvftools/backend/common/log"
	"pvftools/backend/internal/data_source"
)

type WorldDropItem struct {
	Id    int `json:"id"`
	Level int `json:"level"`
	Rate  int `json:"rate"`
}

type worldDropMgr struct {
	parser *parser.Parser
	items  map[int][]*WorldDropItem
}

func (m *worldDropMgr) LoadData() {
	ds := data_source.GetDataSource()
	c, err := ds.GetFileContent(consts.PathWorldDrop)
	if err != nil {
		log.LogError("get file %s content err %v", consts.PathWorldDrop, err)
		return
	}
	m.parser = parser.NewParser(c)
	data, err := m.parser.GetInts(consts.LabelWorldDrop)
	if err != nil {
		log.LogError("parse world drop data err %v", err)
		return
	}
	items := make(map[int][]*WorldDropItem)
	// 每组第一位为 level, 第二位为 0, 之后两两一组为 id-rate, -1 结尾
	levelData := make([][]int, 0)
	tmp := make([]int, 0)
	for _, v := range data {
		if v == -1 {
			levelData = append(levelData, tmp[:])
			tmp = make([]int, 0)
		} else {
			tmp = append(tmp, v)
		}
	}
	for _, v := range levelData {
		level := v[0]
		itemsData := v[2:]
		items[level] = make([]*WorldDropItem, 0)
		for i := 0; i < len(itemsData); i += 2 {
			id := itemsData[i]
			rate := itemsData[i+1]
			items[level] = append(items[level], &WorldDropItem{
				Id:    id,
				Level: level,
				Rate:  rate,
			})
		}
	}
	m.items = items
}

func (m *worldDropMgr) renderWorldDropData() []int {
	data := make([]int, 0)
	for level := 1; level <= 200; level++ {
		data = append(data, level)
		for _, item := range m.GetItems(level) {
			data = append(data, item.Id, item.Rate)
		}
		data = append(data, -1)
	}
	return data
}

func (m *worldDropMgr) SaveData() error {
	m.parser.SetInts(consts.LabelWorldDrop, m.renderWorldDropData())
	return data_source.GetDataSource().SaveFileContent(consts.PathWorldDrop, m.parser.Render())
}

func (m *worldDropMgr) getItemsMap() map[int][]*WorldDropItem {
	if m.items == nil {
		m.LoadData()
	}
	return m.items
}

func (m *worldDropMgr) GetItems(level int) []*WorldDropItem {
	return m.getItemsMap()[level]
}

func (m *worldDropMgr) AddItem(minLev, maxLev, code, rate int) {
	for i := minLev; i <= maxLev; i++ {
		item := &WorldDropItem{
			Id:    code,
			Level: i,
			Rate:  rate,
		}
		m.items[i] = append(m.items[i], item)
	}
}

func (m *worldDropMgr) DelDropItem(minLev, maxLev int, code []int) {
	if maxLev == 0 {
		maxLev = 200
	}
	codeMap := lo.SliceToMap(code, func(v int) (int, bool) { return v, true })
	for i := minLev; i <= maxLev; i++ {
		m.items[i] = lo.Filter(m.items[i], func(item *WorldDropItem, _ int) bool {
			return !codeMap[item.Id]
		})
	}
}

var WorldDropMgr = &worldDropMgr{}
