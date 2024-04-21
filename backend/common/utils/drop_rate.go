package utils

import "github.com/samber/lo"

// CalcPercentageDropRate 计算万分比掉率
func CalcPercentageDropRate(src []int) []int {
	return []int{
		src[0] * 100 / 1000000,
		(src[1] - src[0]) * 100 / 1000000,
		(src[2] - src[1]) * 100 / 1000000,
		(src[3] - src[2]) * 100 / 1000000,
		(src[4] - src[3]) * 100 / 1000000,
	}
}

// CalcPvfDropData 计算pvf掉率
func CalcPvfDropData(src []int) []int {
	return lo.Map([]int{
		src[0],
		src[0] + src[1],
		src[0] + src[1] + src[2],
		src[0] + src[1] + src[2] + src[3],
	}, func(item int, _ int) int { return item * 10000 })
}
