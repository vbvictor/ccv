package level2

func Func5(x int) int {
	a := 0
	for i := 0; i < x; i++ {
		if i > 0 {
			for j := 0; j < i; j++ {
				if j%2 == 0 {
					if j%3 == 0 {
						a += j
					}
				}
			}
		}
	}
	return a
}

func Func6(x int) int {
	a := 0
	for i := 0; i < x; i++ {
		switch i % 3 {
		case 0:
			switch i % 2 {
			case 0:
				a += 1
			default:
				a += 2
			}
		default:
			a += 3
		}
	}
	return a
}
