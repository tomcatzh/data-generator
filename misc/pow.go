package misc

// Pow returns the pow of x an y in integer
func Pow(x, y int) (result int) {
	result = 1

	for i := 0; i < y; i++ {
		result *= x
	}

	return
}
