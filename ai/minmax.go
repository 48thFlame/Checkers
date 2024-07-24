package ai

import (
	"github.com/48thFlame/Checkers/checkers"
)

const (
	piecesCapturedForEndGame = 24 - 8
)

func minMax(g checkers.Game, startDepth, currentDepth int, alpha, beta int) (eval int) {
	MinMaxStatsMan.calls++
	if g.State != checkers.Playing {
		return gameOverEval(g, startDepth, currentDepth)
	}

	if currentDepth == 0 {
		if g.CanCapture() {
			// should not end calculation here - in middle of capture sequence
			MinMaxStatsMan.extendedSearch++
			return minMax(g, startDepth+1, 1, alpha, beta)
		}

		// if not in middle of capture sequence eval
		if g.NPiecesCaptured >= piecesCapturedForEndGame {
			MinMaxStatsMan.endEval++
			return EvaluateEndGamePos(g)
		} else {
			MinMaxStatsMan.midEval++
			return EvaluateMidPosition(g)
		}
	}

	legalMoves := g.GetLegalMoves()

	if g.PlrTurn == checkers.BluePlayer {
		eval = lowestE

		for _, move := range legalMoves {
			futureGame := g
			(&futureGame).PlayMove(move)
			currentMoveEval := minMax(futureGame, startDepth, currentDepth-1, alpha, beta)

			eval = max(currentMoveEval, eval)

			alpha = max(alpha, currentMoveEval)
			if beta <= alpha {
				// Red had a better option in previous branches
				MinMaxStatsMan.alphaBetaBreak++
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
			currentMoveEval := minMax(futureGame, startDepth, currentDepth-1, alpha, beta)

			eval = min(currentMoveEval, eval)

			beta = min(beta, currentMoveEval)
			if beta <= alpha {
				// Blue had better option in previous branches
				MinMaxStatsMan.alphaBetaBreak++
				break
			}
		}

		return eval
	}
}
