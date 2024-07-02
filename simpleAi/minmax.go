package simpleAi

import (
	"github.com/48thFlame/Checkers/checkers"
)

func minMax(g checkers.Game, depth uint, alpha, beta int) (eval int) {
	if g.State != checkers.Playing {
		return gameOverEval(g)
	}

	if depth == 0 {
		return EvaluatePosition(g)
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
