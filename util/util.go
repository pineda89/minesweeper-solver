package util

func InBetween[T uint32 | int | int32 | float32 | float64](r T, i T, i2 T) bool {
	return (r >= i && r <= i2) || (r >= i2 && r <= i)
}
