package level2

func Func1(x int) int {
	a := 0
	if x > 0 {
		if x > 10 {
			a = 1
		} else {
			a = 2
		}
	} else {
		if x < -10 {
			a = 3
		}
	}
	return a
}

func Func2(x int) int {
	a := 0
	for i := 0; i < x; i++ {
		switch i {
		case 1:
			a += 1
		case 2:
			a += 2
		default:
			a += 3
		}
	}
	return a
}
