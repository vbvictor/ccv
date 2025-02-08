package level2

func NestedLoopsWithConditions(n int) int {
	result := 0
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i%2 == 0 {
				if j%2 == 0 {
					result += i * j
				}
			}
		}
	}
	return result
}

func SwitchWithLoops(x int) int {
	sum := 0
	for i := 0; i < x; i++ {
		switch i % 3 {
		case 0:
			if i%2 == 0 {
				sum += i
			}
		case 1:
			if i%4 == 0 {
				sum += i * 2
			}
		default:
			sum += 1
		}
	}
	return sum
}
