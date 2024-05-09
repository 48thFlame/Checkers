package main

func isIn[T comparable](a T, s ...T) bool {
	for _, v := range s {
		if v == a {
			return true
		}
	}

	return false
}
