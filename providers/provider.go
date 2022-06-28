package providers

import "image"

type Provider interface {
	ProcessImage(img image.Image) [][]int
	MoveToCeld(x, y int)
}

const (
	CELL_DISCOVERED  = 0
	CELL_UNKNOWN     = 9
	CELL_1_NEIGHBORN = 1
	CELL_2_NEIGHBORN = 2
	CELL_3_NEIGHBORN = 3
	CELL_4_NEIGHBORN = 4
	CELL_5_NEIGHBORN = 5
	CELL_6_NEIGHBORN = 6
	CELL_FLAG        = 8
)

func GetTag(counters []int) int {
	maxI, _ := GetMax(counters)
	if maxI == CELL_DISCOVERED || maxI == CELL_UNKNOWN {
		newMaxI := 0
		newMaxVal := 0
		for i := 0; i < len(counters); i++ {
			if counters[i] > newMaxVal && i != maxI {
				newMaxI = i
				newMaxVal = counters[i]
			}
		}
		if newMaxI > 0 {
			maxI = newMaxI
		}
	}
	return maxI
}

func GetMax(counters []int) (int, int) {
	var maxI, maxValue = 0, counters[0]
	for i, v := range counters {
		if v > maxValue {
			maxValue = v
			maxI = i
		}
	}
	return maxI, maxValue
}

func GetColor(img image.Image, x int, y int) (uint32, uint32, uint32) {
	color := img.At(x, y)
	rm, gm, bm, _ := color.RGBA()
	return rm / 255, gm / 255, bm / 255
}
