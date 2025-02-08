package level2

func Func7(x, y int) int {
	a := 0
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			switch {
			case i > j:
				if i%2 == 0 {
					a += 1
				}
			case i < j:
				if j%2 == 0 {
					a += 2
				}
			default:
				a += 3
			}
		}
	}
	return a
}

func Func8(x int) int {
	a := 0
	for i := 0; i < x; i++ {
		if i%2 == 0 {
			for j := 0; j < i; j++ {
				switch j {
				case 1:
					if i > 5 {
						if j < 3 {
							a += 1
						}
					}
				default:
					a += 2
				}
			}
		}
	}
	return a
}
