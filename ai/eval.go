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

	nOfMiddleStage = 18
	nOfEndStage    = 8
)

func gameOverEval(g checkers.Game) int {
	switch g.State {
	case checkers.Draw:
		return drawE
	case checkers.BlueWon:
		return blueWonE
	case checkers.RedWon:
		return redWonE
	}

	// * should not get here
	return 0
}

type heatMap [checkers.BoardSize]int

var (
	PiecesHeatMap = heatMap{
		-0, 7, -0, 7, -0, 7, -0, 7,
		3, -0, 1, -0, 1, -0, 1, -0,
		-0, 4, -0, 4, -0, 4, -0, 6,
		10, -0, 12, -0, 12, -0, 10, -0,
		-0, 6, -0, 8, -0, 8, -0, 12,
		5, -0, 0, -0, 0, -0, 0, -0,
		-0, 0, -0, 0, -0, 0, -0, 5,
		0, -0, 0, -0, 0, -0, 0, -0,
	}

	// MiddleHeatMap = heatMap{
	// 	-0, 3, -0, 3, -0, 3, -0, 3,
	// 	0, -0, 0, -0, 0, -0, 0, -0,
	// 	-0, 2, -0, 2, -0, 2, -0, 2,
	// 	11, -0, 10, -0, 10, -0, 8, -0,
	// 	-0, 8, -0, 10, -0, 10, -0, 11,
	// 	7, -0, 4, -0, 4, -0, 4, -0,
	// 	-0, 0, -0, 0, -0, 0, -0, 1,
	// 	0, -0, 0, -0, 0, -0, 0, -0,
	// }

	KingHeatMap = heatMap{
		-0, 1, -0, -10, -0, -11, -0, -12,
		1, -0, 4, -0, 3, -0, 0, -0,
		-0, 4, -0, 8, -0, 8, -0, -11,
		-10, -0, 12, -0, 12, -0, 7, -0,
		-0, 7, -0, 12, -0, 12, -0, -10,
		-11, -0, 8, -0, 8, -0, 4, -0,
		-0, 0, -0, 3, -0, 4, -0, 1,
		-12, -0, -11, -0, -10, -0, 1, -0,
	}
)

func EvaluatePosition(g checkers.Game) (eval int) {
	// var nBlue, nRed uint8

	// blues := make([]int, 0)
	// reds := make([]int, 0)

	for slotI, slot := range g.Board {
		switch slot {
		case checkers.BluePiece:
			// blues = append(blues, slotI)

			// nBlue++
			eval += pieceWeightE
			eval += PiecesHeatMap[slotI]

		case checkers.BlueKing:
			// blues = append(blues, slotI)

			// nBlue++
			eval += kingWeightE
			eval += KingHeatMap[slotI]

		case checkers.RedPiece:
			// reds = append(reds, slotI)

			// nRed++
			eval -= pieceWeightE
			eval -= PiecesHeatMap[checkers.BoardSize-1-slotI]

		case checkers.RedKing:
			// reds = append(reds, slotI)

			// nRed++
			eval -= kingWeightE
			eval -= KingHeatMap[checkers.BoardSize-1-slotI]

		}
	}

	// totalPiecesN := nBlue + nRed
	// var hm heatMap

	// // if totalPiecesN > nOfMiddleStage {
	// // 	// opening
	// if totalPiecesN > nOfEndStage {
	// 	hm = PiecesHeatMap
	// 	// 	// middle
	// 	// 	// hm = MiddleHeatMap
	// } else {
	// 	// 	// 	// 	// end
	// 	hm = KingHeatMap

	// }

	// for _, blueI := range blues {
	// 	eval += hm[blueI]
	// }

	// for _, redI := range reds {
	// 	eval -= hm[checkers.BoardSize-1-redI]
	// }
	// eval += OpeningHeatMap[slotI]
	// eval += OpeningHeatMap[slotI]
	// eval -= OpeningHeatMap[checkers.BoardSize-1-slotI]
	// eval -= OpeningHeatMap[checkers.BoardSize-1-slotI]
	return eval
}
