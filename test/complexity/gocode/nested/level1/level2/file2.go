package level2

func Func3(x int) int {
	a := 0
	for i := 0; i < x; i++ {
		for j := 0; j < x; j++ {
			if i > j {
				a += 1
			}
		}
	}
	return a
}

func Func4(x int) int {
	a := 0
	for i := 0; i < x; i++ {
		if i%2 == 0 {
			switch i {
			case 0:
				a += 1
			default:
				if i > 5 {
					a += 2
				}
			}
		}
	}
	return a
}
