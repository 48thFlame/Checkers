package ai

import (
	"github.com/48thFlame/Checkers/checkers"
)

const (
	highestE = 199_999_999
	lowestE  = highestE * -1

	blueWonE = 1_000_000
	redWonE  = blueWonE * -1
	drawE    = 0

	pieceWeightE = 100
	kingWeightE  = 150
)

func gameOverEval(g checkers.Game, startDepth, currentDepth int) int {
	switch g.State {
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

type heatMap [checkers.BoardSize]int

var (
	PiecesHeatMap = heatMap{
		-0, 7, -0, 7, -0, 7, -0, 7,
		1, -0, 0, -0, 0, -0, 0, -0,
		-0, 4, -0, 5, -0, 5, -0, 3,
		11, -0, 12, -0, 12, -0, 10, -0,
		-0, 3, -0, 4, -0, 4, -0, 12,
		5, -0, 1, -0, 1, -0, 1, -0,
		-0, 20, -0, 20, -0, 20, -0, 20,
		0, -0, 0, -0, 0, -0, 0, -0,
	}

	KingHeatMap = heatMap{
		-0, -3, -0, -10, -0, -11, -0, -12,
		-3, -0, 4, -0, 3, -0, -1, -0,
		-0, 4, -0, 8, -0, 8, -0, -11,
		-10, -0, 13, -0, 13, -0, 7, -0,
		-0, 7, -0, 13, -0, 13, -0, -10,
		-11, -0, 8, -0, 8, -0, 4, -0,
		-0, -1, -0, 3, -0, 4, -0, -3,
		-12, -0, -11, -0, -10, -0, -3, -0,
	}
)

func EvaluateMidPosition(g checkers.Game) (eval int) {
	for slotI, slot := range g.Board {
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
		9, -0, 1, -0, 0, -0, 0, -0,
		-0, 1, -0, 2, -0, 2, -0, -10,
		-10, -0, 2, -0, 2, -0, 0, -0,
		-0, 0, -0, 2, -0, 2, -0, -10,
		-11, -0, 2, -0, 2, -0, 1, -0,
		-0, 0, -0, 0, -0, 1, -0, 9,
		-12, -0, -11, -0, -10, -0, 9, -0,
	}
)

func EvaluateEndGamePos(g checkers.Game) (eval int) {
	var nBlue, nRed uint

	blueKingIs := make([]int, 0)
	redKingIs := make([]int, 0)

	for slotI, slot := range g.Board {
		switch slot {
		case checkers.BluePiece:
			nBlue++
			eval += endPieceWeightE

		case checkers.RedPiece:
			nRed++
			eval -= endPieceWeightE

		case checkers.BlueKing:
			nBlue++

			eval += endKingWeightE
			eval += EndKingHeatMap[slotI]

			blueKingIs = append(blueKingIs, slotI)

		case checkers.RedKing:
			nRed++

			eval -= endKingWeightE
			eval -= EndKingHeatMap[checkers.BoardSize-1-slotI]

			redKingIs = append(redKingIs, slotI)
		}
	}

	if nBlue != nRed {
		var dist, distScore int

		for _, bki := range blueKingIs {
			for _, rki := range redKingIs {
				dist = getManhattanDist(bki, rki)

				if dist >= 2 {
					distScore += 9 - dist // punish being far, reward close
				}
			}
		}

		if nBlue > nRed {
			eval += distScore
		} else { // => nBlue < nRed
			eval -= distScore
		}
	}

	return eval
}
