package ai

import (
	"github.com/48thFlame/Checkers/checkers"
)

//TODO: https://g.co/gemini/share/d5b2955d26b2

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
	piecesHeatMap = heatMap{
		-0, 8, -0, 8, -0, 8, -0, 6,
		2, -0, 2, -0, 2, -0, 1, -0,
		-0, 6, -0, 6, -0, 6, -0, 6,
		9, -0, 10, -0, 10, -0, 10, -0,
		-0, 4, -0, 5, -0, 5, -0, 11,
		3, -0, 0, -0, 0, -0, 0, -0,
		-0, 16, -0, 18, -0, 18, -0, 15, // hoping doesn't end calculation when can be captured
		0, -0, 0, -0, 0, -0, 0, -0,
	}

	kingHeatMap = heatMap{
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

func evaluateMidPosition(agd aiGameData) (eval int) {
	for slotI, slot := range agd.g.Board {
		switch slot {
		case checkers.BluePiece:
			eval += piecesHeatMap[slotI]
			eval += pieceWeightE

		case checkers.BlueKing:
			eval += kingWeightE
			eval += kingHeatMap[slotI]

		case checkers.RedPiece:
			eval -= pieceWeightE
			eval -= piecesHeatMap[checkers.BoardSize-1-slotI]

		case checkers.RedKing:
			eval -= kingWeightE
			eval -= kingHeatMap[checkers.BoardSize-1-slotI]
		}
	}

	return eval
}

const (
	endPieceWeightE = 70
	endKingWeightE  = 150
	endPiecePunishE = 5
)

var (
	endKingHeatMap = heatMap{
		-0, 3, -0, -10, -0, -11, -0, -12,
		3, -0, 1, -0, 0, -0, -2, -0,
		-0, 1, -0, 2, -0, 2, -0, -10,
		-10, -0, 2, -0, 2, -0, 0, -0,
		-0, 0, -0, 2, -0, 2, -0, -10,
		-11, -0, 2, -0, 2, -0, 1, -0,
		-0, -2, -0, 0, -0, 1, -0, 3,
		-12, -0, -11, -0, -10, -0, 3, -0,
	}
)

func iAbs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func getManhattanDist(a, b int) int {
	aCol := a % checkers.BoardSideSize
	aRow := a / checkers.BoardSideSize

	bCol := b % checkers.BoardSideSize
	bRow := b / checkers.BoardSideSize

	deltaCol := iAbs(aCol - bCol)
	deltaRow := iAbs(aRow - bRow)

	return deltaCol + deltaRow
}

func evaluateEndGamePos(agd aiGameData) (eval int) {
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
			eval += endKingHeatMap[slotI]

			blueKingIs = append(blueKingIs, slotI)

		case checkers.RedKing:
			eval -= endKingWeightE
			eval -= endKingHeatMap[checkers.BoardSize-1-slotI]

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
			eval -= agd.nRed * endPiecePunishE // reward trading
		} else { // nBlue < nRed
			eval -= distScore
			eval += agd.nBlue * endPiecePunishE
		}
	}

	return eval
}
