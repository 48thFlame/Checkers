package ai

import (
	"github.com/48thFlame/Checkers/checkers"
)

// classic min-max with alpha beta pruning
func minMax(agd aiGameData, startingLegalMoves []checkers.Move, startDepth, currentDepth int, alpha, beta int) (me moveEval) {
	MinMaxStatsMan.calls++

	if agd.g.State != checkers.Playing {
		MinMaxStatsMan.gameEval++
		me.eval = gameOverEval(agd, startDepth, currentDepth)
		return me
	}

	if currentDepth == 0 {
		if agd.canCapture() {
			// should not end calculation here - in middle of capture sequence
			// extend search for 1 more move
			MinMaxStatsMan.extendedSearch++
			return minMax(agd, startingLegalMoves, startDepth+1, 1, alpha, beta)
		}
		// if not in middle of capture sequence regular eval

		if agd.isInEndGame() {
			MinMaxStatsMan.endEval++
			me.eval = EvaluateEndGamePos(agd)
			return me
		} else {
			MinMaxStatsMan.midEval++
			me.eval = EvaluateMidPosition(agd)
			return me
		}
	}

	var legalMoves []checkers.Move
	if startDepth == currentDepth { // in starting position - use those moves there optimized with best move first
		legalMoves = startingLegalMoves
	} else {
		legalMoves = agd.g.GetLegalMoves()
	}

	if agd.g.PlrTurn == checkers.BluePlayer {
		me.eval = lowestE

		for _, move := range legalMoves {
			futureGame := agd
			(&futureGame).playMove(move)
			currentMoveEval := minMax(futureGame, startingLegalMoves, startDepth, currentDepth-1, alpha, beta)

			if currentMoveEval.eval > me.eval { // if current move is better then previously checked
				me.move = move // its new best move
				me.eval = currentMoveEval.eval
			}

			alpha = max(alpha, currentMoveEval.eval)
			if beta <= alpha {
				MinMaxStatsMan.alphaBetaBreak++
				break
			}

		}

		me.depth = startDepth
		return me
	} else {
		// reds minimizing turn
		me.eval = highestE

		for _, move := range legalMoves {
			futureGame := agd
			(&futureGame).playMove(move)
			currentMoveEval := minMax(futureGame, startingLegalMoves, startDepth, currentDepth-1, alpha, beta)

			if currentMoveEval.eval < me.eval {
				me.move = move
				me.eval = currentMoveEval.eval
			}

			beta = min(beta, currentMoveEval.eval)
			if beta <= alpha {
				MinMaxStatsMan.alphaBetaBreak++
				break
			}
		}

		me.depth = startDepth
		return me
	}
}
