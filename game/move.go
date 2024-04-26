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
	upLeftCalc  = -BoardSideSize - 1
	upRightCalc = -BoardSideSize + 1

	downLeftCalc  = BoardSideSize - 1
	downRightCalc = BoardSideSize + 1
)

// these are needed to know when should become king
var (
	// TODO: what if board size is not 8?!!
	topBoardEndI    = [...]int{1, 3, 5, 7}
	bottomBoardEndI = [...]int{56, 58, 60, 62}
)

func newMove(startI, endI int, capturedI ...int) Move {
	return Move{startI: startI, endI: endI, capturedPiecesI: capturedI}
}

type Move struct {
	startI int
	endI   int

	capturedPiecesI []int
}

func getLegalDownMovingMoves(b Board, i int) []Move {
	moves := make([]Move, 0)

	downLeftI := i + downLeftCalc
	downRightI := i + downRightCalc

	if downLeftI < BoardSize && b[downLeftI] == Empty {
		moves = append(moves, newMove(i, downLeftI))
	}

	if downRightI < BoardSize && b[downRightI] == Empty {
		moves = append(moves, newMove(i, downRightI))
	}

	return moves
}

func getLegalUpMovingMoves(b Board, i int) []Move {
	moves := make([]Move, 0)

	upLeftI := i + upLeftCalc
	upRightI := i + upRightCalc

	if upLeftI > 0 && b[upLeftI] == Empty {
		moves = append(moves, newMove(i, upLeftI))
	}

	if upRightI > 0 && b[upRightI] == Empty {
		moves = append(moves, newMove(i, upRightI))
	}

	return moves
}

// GetLegalMoves returns a slice all legal moves in position
func (g Game) GetLegalMoves() []Move {
	moves := make([]Move, 0)

	// can start at 1, because 0 (top left) is always a NaS, same reason stops 1 early
	for i := 1; i < BoardSize-1; i++ {
		slot := g.Board[i]

		if slot == NaS || slot == Empty {
			continue
		}

		switch g.PlrTurn {
		case BluePlayer:
			// blue goes down
			if isIn(slot, BluePiece, BlueKing) {
				moves = append(moves, getLegalDownMovingMoves(g.Board, i)...)

				// if its specifically a king then also check up
				if slot == BlueKing {
					moves = append(moves, getLegalUpMovingMoves(g.Board, i)...)
				}
			}

		case RedPlayer:
			// red goes up
			if isIn(slot, RedPiece, RedKing) {
				moves = append(moves, getLegalUpMovingMoves(g.Board, i)...)

				// if its specifically a king then also check down
				if slot == RedKing {
					moves = append(moves, getLegalDownMovingMoves(g.Board, i)...)
				}
			}
		}

	}

	return moves
}
