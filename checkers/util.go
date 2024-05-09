package checkers

func isIn[T comparable](a T, s ...T) bool {
	for _, v := range s {
		if v == a {
			return true
		}
	}

	return false
}

// TODO: make better
func isOnEnd(plr Player, i int) bool {
	switch plr {
	case BluePlayer:
		return isIn(i, bottomBoardEndI[:]...)

	case RedPlayer:
		return isIn(i, topBoardEndI[:]...)
	}

	return false
}
