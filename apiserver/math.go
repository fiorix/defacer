package apiserver

func roundUp(v int) int {
	return 10 * ((v + 9) / 10)
}

func roundDown(v int) int {
	return 10 * (v / 10)
}
