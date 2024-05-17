package ai

import (
	"github.com/48thFlame/Checkers/checkers"
)

const (
	highestE = 100_000
	lowestE  = highestE * -1

	blueWonE = 10_000
	redWonE  = blueWonE * -1
	drawE    = 0

	pieceWeightE = 10
	kingWeightE  = 15
)

type pieceCount struct {
	bPieces, bKings int
	rPieces, rKings int
}

func evaluatePosition(g checkers.Game) int {
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

func minMax(g checkers.Game, depth uint, alpha, beta int) (eval int) {
	if g.State != checkers.Playing {
		return gameOverEval(g)
	}

	if depth == 0 {
		return evaluatePosition(g)
	}

	legalMoves := g.GetLegalMoves()

	if g.PlrTurn == checkers.BluePlayer {
		eval = lowestE

		for _, move := range legalMoves {
			futureGame := g
			(&futureGame).PlayMove(move)
			currentMoveEval := minMax(futureGame, depth-1, alpha, beta)

			eval = max(currentMoveEval, eval)

			alpha = max(alpha, currentMoveEval)
			if beta <= alpha {
				// Red had a better option in previous branches
				break
			}
		}

		return eval
	} else {
		// reds minimizing turn
		eval = highestE

		for _, move := range legalMoves {
			futureGame := g
			(&futureGame).PlayMove(move)
			currentMoveEval := minMax(futureGame, depth-1, alpha, beta)

			eval = min(currentMoveEval, eval)

			beta = min(beta, currentMoveEval)
			if beta <= alpha {
				// Blue had better option in previous branches
				break
			}
		}

		return eval
	}
}
