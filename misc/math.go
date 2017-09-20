package misc

// MinInt returns the minimal int value
func MinInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}
