package game

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
		if isIn(i, bottomBoardEndI[:]...) {
			return true
		}
	case RedPlayer:
		if isIn(i, topBoardEndI[:]...) {
			return true
		}
	}

	return false
}
