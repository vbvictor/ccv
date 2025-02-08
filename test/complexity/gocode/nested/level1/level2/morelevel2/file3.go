package level2

func ComplexNestedStructure(x, y int) int {
	result := 0
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			if i%2 == 0 {
				if j%2 == 0 {
					switch i % 3 {
					case 0:
						result += i * j
					case 1:
						if j%3 == 0 {
							result += i + j
						}
					}
				}
			}
		}
	}
	return result
}

func MultipleControlFlows(n int) int {
	result := 0
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			for j := 0; j < i; j++ {
				if j%3 == 0 {
					switch j % 4 {
					case 0:
						result += 1
					case 1:
						if i%5 == 0 {
							result += 2
						}
					default:
						result += 3
					}
				}
			}
		}
	}
	return result
}
