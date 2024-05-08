package game

/*
Safe to check also edges because if its really an edge and cant go to a side,
after calculation its ends up in a Nas.

cn = colsNum = BoardSideSize

-cn-1 -cn -cn+1
-1 0 1
cn-1 cn cn+1
*/
const (
	downLeftCalc  = BoardSideSize - 1
	downRightCalc = BoardSideSize + 1

	upLeftCalc  = -BoardSideSize - 1
	upRightCalc = -BoardSideSize + 1
)

var (
	// blue goes down
	blueDirCalcs = [...]int{downLeftCalc, downRightCalc}
	// red goes up
	redDirCalcs = [...]int{upLeftCalc, upRightCalc}

	kingDirCalcs = [...]int{downLeftCalc, downRightCalc, upLeftCalc, downRightCalc}
)

// these are needed to know when should become king
var (
	// TODO: what if board size is not 8?!!
	topBoardEndI    = [...]int{1, 3, 5, 7}
	bottomBoardEndI = [...]int{56, 58, 60, 62}
)

var (
	redPieces  = [...]BoardSlot{RedPiece, RedKing}
	bluePieces = [...]BoardSlot{BluePiece, BlueKing}
)

func newMove(startI, endI int, capturedI ...int) Move {
	return Move{startI: startI, endI: endI, capturedPiecesI: capturedI}
}

type Move struct {
	startI int
	endI   int

	capturedPiecesI []int
}

// returns a slot and its I (after calc)
func getSlotAfterCalc(b Board, i, dirCalc int) (BoardSlot, int) {
	newI := i + dirCalc

	// 0 is always a Nas so < and not <=
	if 0 < newI && newI < BoardSize {
		return b[newI], newI
	} else {
		return NaS, -1
	}
}

// returns `directionsToUse` and whether slot matches plr - if not wrong turn and wrong slot
func getDirectionsToUse(plr Player, slot BoardSlot) (directionsToUse []int, match bool) {
	switch plr {
	case BluePlayer:
		if slot == BluePiece {
			directionsToUse = blueDirCalcs[:]
			match = true
		} else if slot == BlueKing {
			directionsToUse = kingDirCalcs[:]
			match = true
		}
	case RedPlayer:
		if slot == RedPiece {
			directionsToUse = redDirCalcs[:]
			match = true
		} else if slot == RedKing {
			directionsToUse = kingDirCalcs[:]
			match = true
		}
	}

	return
}

// get legal *moving* moves
func getMovings(b Board, i int, dirCalcs []int) []Move {
	moves := make([]Move, 0)

	for _, dc := range dirCalcs {
		dSlot, dI := getSlotAfterCalc(b, i, dc)

		if dSlot == Empty {
			moves = append(moves, newMove(i, dI))
		}
	}

	return moves
}

// get legal capturing moves
func getCaptures(b Board, i int, directionCalcs []int, enemyPieces []BoardSlot) []Move {
	captures := make([]Move, 0)

	for _, dirCalc := range directionCalcs {
		eatSlot, eatI := getSlotAfterCalc(b, i, dirCalc)
		landSlot, landI := getSlotAfterCalc(b, i, dirCalc*2)

		if isIn(eatSlot, enemyPieces...) && landSlot == Empty { // if can eat
			// clear board so wont eat again
			b[landI] = b[i]
			b[i] = Empty
			b[eatI] = Empty

			// check for new captures in new position, then join em all
			secondLevelCaptures := getCaptures(b, landI, directionCalcs, enemyPieces)

			if len(secondLevelCaptures) > 0 {
				// join em all
				for _, slc := range secondLevelCaptures {
					captures = append(captures,
						newMove(i, slc.endI, append(slc.capturedPiecesI, eatI)...))
				}
			} else {
				captures = append(captures, newMove(i, landI, eatI))
			}
		}
	}
	return captures
}

// GetLegalMoves returns a slice of all legal moves in position
func (g Game) GetLegalMoves() []Move {
	moves := make([]Move, 0)

	// can the current player capture
	canCapture := false

	// can start at 1, because 0 (top left) is always a NaS, same reason stops 1 early
	for i := 1; i < BoardSize-1; i++ {
		slot := g.Board[i]

		if slot == NaS || slot == Empty {
			continue
		}

		var directionsToUse []int
		var good bool // is looking maybe at wrong slot because not thats player turn?
		var enemyPieces []BoardSlot

		switch g.PlrTurn {
		case BluePlayer:
			directionsToUse, good = getDirectionsToUse(BluePlayer, slot)
			if !good {
				continue
			}

			enemyPieces = redPieces[:]
		case RedPlayer:
			directionsToUse, good = getDirectionsToUse(RedPlayer, slot)
			if !good {
				continue
			}

			enemyPieces = bluePieces[:]
		}

		captures := getCaptures(g.Board, i, directionsToUse, enemyPieces)
		capturesLen := len(captures)
		if capturesLen > 0 {
			if !canCapture {
				// if its first capture - clear all moves - they are no longer legal!
				moves = make([]Move, 0, capturesLen)
				canCapture = true
			}

			moves = append(moves, captures...)
			continue
		}

		if !canCapture {
			// because if can capture can't move
			moves = append(moves, getMovings(g.Board, i, directionsToUse)...)
		}
	}

	return moves
}
