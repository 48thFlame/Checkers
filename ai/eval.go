package ai

import (
	"fmt"

	"github.com/48thFlame/Checkers/checkers"
)

type moveEval struct {
	depth int
	move  checkers.Move
	eval  int
}

func (me moveEval) String() string {
	return fmt.Sprintf("(%d| %d,%d |%d)",
		me.depth, me.move.StartI, me.move.EndI, me.eval)
}

const (
	highestE = 199_999_999
	lowestE  = highestE * -1
)

const (
	blueWonE = 1_000_000
	redWonE  = blueWonE * -1
	drawE    = 0
)

func gameOverEval(agd aiGameData, startDepth, currentDepth int) int {
	switch agd.g.State {
	case checkers.Draw:
		return drawE
	case checkers.BlueWon:
		return blueWonE - startDepth + currentDepth // this way the fastest win is the best win
	case checkers.RedWon:
		return redWonE + startDepth - currentDepth
	}

	// * should not get here
	return 0
}

const (
	pieceWeightE = 90
	kingWeightE  = 150
)

type heatMap [checkers.BoardSize]int

var (
	PiecesHeatMap = heatMap{
		-0, 8, -0, 8, -0, 8, -0, 6,
		2, -0, 2, -0, 2, -0, 1, -0,
		-0, 6, -0, 6, -0, 6, -0, 6,
		9, -0, 10, -0, 10, -0, 10, -0,
		-0, 4, -0, 5, -0, 5, -0, 11,
		3, -0, 0, -0, 0, -0, 0, -0,
		-0, 25, -0, 30, -0, 30, -0, 20, // hoping doesn't end calculation when can be captured
		0, -0, 0, -0, 0, -0, 0, -0,
	}

	KingHeatMap = heatMap{
		-0, -4, -0, -5, -0, -5, -0, -5,
		-4, -0, 2, -0, 2, -0, -1, -0,
		-0, 3, -0, 5, -0, 5, -0, -3,
		-3, -0, 6, -0, 6, -0, 4, -0,
		-0, 4, -0, 6, -0, 6, -0, -3,
		-3, -0, 5, -0, 5, -0, 3, -0,
		-0, -1, -0, 2, -0, 2, -0, -4,
		-5, -0, -5, -0, -5, -0, -4, -0,
	}
)

func EvaluateMidPosition(agd aiGameData) (eval int) {
	for slotI, slot := range agd.g.Board {
		switch slot {
		case checkers.BluePiece:
			eval += PiecesHeatMap[slotI]
			eval += pieceWeightE

		case checkers.BlueKing:
			eval += kingWeightE
			eval += KingHeatMap[slotI]

		case checkers.RedPiece:
			eval -= pieceWeightE
			eval -= PiecesHeatMap[checkers.BoardSize-1-slotI]

		case checkers.RedKing:
			eval -= kingWeightE
			eval -= KingHeatMap[checkers.BoardSize-1-slotI]
		}
	}

	return eval
}

const (
	endPieceWeightE = 70
	endKingWeightE  = 150
)

var (
	EndKingHeatMap = heatMap{
		-0, 9, -0, -10, -0, -11, -0, -12,
		9, -0, 1, -0, 0, -0, -2, -0,
		-0, 1, -0, 2, -0, 2, -0, -10,
		-10, -0, 2, -0, 2, -0, 0, -0,
		-0, 0, -0, 2, -0, 2, -0, -10,
		-11, -0, 2, -0, 2, -0, 1, -0,
		-0, -2, -0, 0, -0, 1, -0, 9,
		-12, -0, -11, -0, -10, -0, 9, -0,
	}
)

func EvaluateEndGamePos(agd aiGameData) (eval int) {
	blueKingIs := make([]int, 0)
	redKingIs := make([]int, 0)

	for slotI, slot := range agd.g.Board {
		switch slot {
		case checkers.BluePiece:
			eval += endPieceWeightE

		case checkers.RedPiece:
			eval -= endPieceWeightE

		case checkers.BlueKing:
			eval += endKingWeightE
			eval += EndKingHeatMap[slotI]

			blueKingIs = append(blueKingIs, slotI)

		case checkers.RedKing:
			eval -= endKingWeightE
			eval -= EndKingHeatMap[checkers.BoardSize-1-slotI]

			redKingIs = append(redKingIs, slotI)
		}
	}

	if agd.nBlue != agd.nRed {
		var dist, distScore int

		for _, bki := range blueKingIs {
			for _, rki := range redKingIs {
				dist = getManhattanDist(bki, rki)

				if dist >= 2 {
					distScore += 7 - dist // punish being far, reward close
				}
			}
		}

		if agd.nBlue > agd.nRed {
			eval += distScore
		} else { // => nBlue < nRed
			eval -= distScore
		}
	}

	return eval
}
