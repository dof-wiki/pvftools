package utils

// CalcPercentageDropRate 计算万分比掉率
func CalcPercentageDropRate(src []int) []int {
	return []int{
		src[0] * 10000 / 1000000,
		(src[1] - src[0]) * 10000 / 1000000,
		(src[2] - src[1]) * 10000 / 1000000,
		(src[3] - src[2]) * 10000 / 1000000,
		(src[4] - src[3]) * 10000 / 1000000,
	}
}

// CalcPvfDropData 计算pvf掉率
func CalcPvfDropData(src []int) []int {
	return []int{
		src[0],
		src[0] + src[1],
		src[0] + src[1] + src[2],
		src[0] + src[1] + src[2] + src[3],
	}
}
