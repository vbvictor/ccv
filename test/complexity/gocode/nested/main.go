package nested

func BaseFunction() int {
	return 42
}

func SimpleCondition(x int) int {
	if x > 0 {
		return x
	}
	return -x
}
