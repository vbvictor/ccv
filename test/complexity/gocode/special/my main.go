package special

func simpleFunction(x int) int {
	return x + 1
}

func complexFunction(n int) int {
	if n <= 1 {
		return n
	}

	result := 0
	for i := 1; i <= n; i++ {
		if i%2 == 0 {
			if i%3 == 0 {
				result += i * 2
			} else {
				result += i
			}
		} else if i%3 == 0 {
			result += i * 3
		} else {
			switch {
			case i%5 == 0:
				result += i * 5
			case i%7 == 0:
				result += i * 7
			default:
				result += i
			}
		}
	}
	return result
}
