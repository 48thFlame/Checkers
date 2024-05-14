package ai

import (
	"math"

	"github.com/48thFlame/Checkers/checkers"
)

const (
	highestE = 1000
	lowestE  = highestE * -1

	blueWonE = 100
	redWonE  = blueWonE * -1
	drawE    = 0
)

type pieceCount struct {
	bp, bk float64
	rb, rk float64
}

func evaluatePosition(g checkers.Game) float64 {
	pc := pieceCount{}

	for _, slot := range g.Board {
		switch slot {
		case checkers.BluePiece:
			pc.bp++
		case checkers.BlueKing:
			pc.bk++
		case checkers.RedPiece:
			pc.rb++
		case checkers.RedKing:
			pc.rk++
		}
	}

	return pc.bp + (pc.bk * 1.5) - pc.rb - (pc.rk * 1.5)
}

func gameOverEval(g checkers.Game) float64 {
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

func minMax(g checkers.Game, depth uint, alpha, beta float64) (eval float64) {
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

			eval = math.Max(currentMoveEval, eval)

			alpha = math.Max(alpha, currentMoveEval)
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

			eval = math.Min(currentMoveEval, eval)

			beta = math.Min(beta, currentMoveEval)
			if beta <= alpha {
				// Blue had better option in previous branches
				break
			}
		}

		return eval
	}
}
