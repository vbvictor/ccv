package level1

func NestedIf(x int) int {
	if x > 10 {
		if x > 20 {
			return 2
		}
		return 1
	}
	return 0
}

func LoopWithCondition(n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			if i%3 == 0 {
				sum += i
			}
		}
	}
	return sum
}
