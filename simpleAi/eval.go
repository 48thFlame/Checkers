package simpleAi

import "github.com/48thFlame/Checkers/checkers"

const (
	pieceWeightE = 100
	kingWeightE  = 150
)

const (
	highestE = 100_999_999
	lowestE  = highestE * -1

	blueWonE = 100_000
	redWonE  = blueWonE * -1
	drawE    = 0
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

type pieceCount struct {
	bPieces, bKings int
	rPieces, rKings int
}

func EvaluatePosition(g checkers.Game) int {
	pc := pieceCount{}

	for _, slot := range g.Board {
		switch slot {
		case checkers.BluePiece:
			pc.bPieces++
		case checkers.BlueKing:
			pc.bKings++
		case checkers.RedPiece:
			pc.rPieces++
		case checkers.RedKing:
			pc.rKings++
		}
	}

	return (pc.bPieces * pieceWeightE) + (pc.bKings * kingWeightE) -
		(pc.rPieces * pieceWeightE) - (pc.rKings * kingWeightE)
}
