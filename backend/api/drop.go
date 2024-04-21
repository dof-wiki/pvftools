package api

import (
	"errors"
	"pvftools/backend/common/consts"
	"pvftools/backend/common/log"
	"pvftools/backend/common/utils"
	"pvftools/backend/internal/world_drop"
)

func (a *App) GetHellDrop() ([]int, error) {
	p := a.getParser(consts.PathHellDrop)
	data, err := p.GetInts(consts.LabelBasisOfRarityDicision)
	if err != nil {
		log.LogError("数据错误%v", err)
		return nil, err
	}
	if len(data) != 13 {
		log.LogError("数据错误")
		return nil, errors.New("数据错误")
	}
	hard := utils.CalcPercentageDropRate(data[1:6])
	veryHard := utils.CalcPercentageDropRate(data[7:12])
	return append(hard, veryHard...), nil
}

func (a *App) SetHellDrop(hard, veryHard []int) error {
	p := a.getParser(consts.PathHellDrop)
	val := []int{2}
	for _, v := range utils.CalcPvfDropData(hard) {
		val = append(val, v)
	}
	val = append(val, 1000000, 1000001)
	for _, v := range utils.CalcPvfDropData(veryHard) {
		val = append(val, v)
	}
	val = append(val, 1000001, 1000002)
	p.SetInts(consts.LabelBasisOfRarityDicision, val)
	return a.saveData(consts.PathHellDrop)
}

func (a *App) GetMonsterDrop() ([][]int, error) {
	p := a.getParser(consts.PathMonsterDrop)
	data, err := p.GetInts(consts.LabelBasisOfRarityDicision)
	if err != nil {
		log.LogError("数据加载失败%v", err)
		return nil, err
	}
	rates := make([][]int, 0, 4)
	for i := 0; i < 4; i++ {
		rates = append(rates, utils.CalcPercentageDropRate(data[i*6:(i+1)*6]))
	}
	return rates, nil
}

func (a *App) SetMonsterDrop(rates [][]int) error {
	p := a.getParser(consts.PathMonsterDrop)
	data := make([]int, 0, 24)
	for _, rate := range rates {
		pvfRate := utils.CalcPvfDropData(rate)
		data = append(data, pvfRate...)
		data = append(data, 1000001, 1000002)
	}
	p.SetInts(consts.LabelBasisOfRarityDicision, data)
	return a.saveData(consts.PathMonsterDrop)
}

func (a *App) GetExtraMonsterDrop() ([][]int, error) {
	p := a.getParser(consts.PathMonsterExtraDrop)
	data, err := p.GetInts(consts.LabelBasisOfRarityDicision)
	if err != nil {
		return nil, err
	}
	rates := make([][]int, 0, 4)
	for i := 0; i < 4; i++ {
		rates = append(rates, utils.CalcPercentageDropRate(data[i*6:(i+1)*6]))
	}
	return rates, nil
}

func (a *App) SetMonsterExtraDrop(rates [][]int) error {
	p := a.getParser(consts.PathMonsterExtraDrop)
	data := make([]int, 0, 24)
	for _, rate := range rates {
		pvfRate := utils.CalcPvfDropData(rate)
		data = append(data, pvfRate...)
		data = append(data, 1000001, 1000002)
	}
	p.SetInts(consts.LabelBasisOfRarityDicision, data)
	return a.saveData(consts.PathMonsterExtraDrop)
}

func (a *App) GetClearRewardDrop() ([]int, error) {
	p := a.getParser(consts.PathClearRewardDrop)
	data, err := p.GetInts(consts.LabelBasisOfRarityDicision)
	if err != nil {
		return nil, err
	}
	return utils.CalcPercentageDropRate(data), nil
}

func (a *App) SetClearRewardDrop(rates []int) error {
	p := a.getParser(consts.PathClearRewardDrop)
	pvfRate := utils.CalcPvfDropData(rates)
	pvfRate = append(pvfRate, 1000001, 1000002)
	p.SetInts(consts.LabelBasisOfRarityDicision, pvfRate)
	return a.saveData(consts.PathClearRewardDrop)
}

func (a *App) GetWorldDropData(level int) ([]*world_drop.WorldDropItem, error) {
	return world_drop.WorldDropMgr.GetItems(level), nil
}

func (a *App) AddWorldDropItem(minLevel, maxLevel, code, rate int) error {
	world_drop.WorldDropMgr.AddItem(minLevel, maxLevel, code, rate)
	return nil
}

func (a *App) DelWorldDropItem(minLevel, maxLevel int, codes []int) error {
	world_drop.WorldDropMgr.DelDropItem(minLevel, maxLevel, codes)
	return nil
}

func (a *App) SaveWorldDropData() error {
	return world_drop.WorldDropMgr.SaveData()
}
