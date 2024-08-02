package ai

import (
	"github.com/48thFlame/Checkers/checkers"
)

const (
	piecesCapturedForEndGame = 24 - 8
)

// classic min-max with alpha beta pruning
func minMax(g checkers.Game, startingLegalMoves []checkers.Move, startDepth, currentDepth int, alpha, beta int) (me moveEval) {
	MinMaxStatsMan.calls++

	if g.State != checkers.Playing {
		MinMaxStatsMan.gameEval++
		me.eval = gameOverEval(g, startDepth, currentDepth)
		return me
	}

	if currentDepth == 0 {
		if g.CanCapture() {
			// should not end calculation here - in middle of capture sequence
			// extend search for 1 more move
			MinMaxStatsMan.extendedSearch++
			return minMax(g, startingLegalMoves, startDepth+1, 1, alpha, beta)
		}

		// if not in middle of capture sequence eval
		if g.NPiecesCaptured >= piecesCapturedForEndGame {
			MinMaxStatsMan.endEval++
			me.eval = EvaluateEndGamePos(g)
			return me
		} else {
			MinMaxStatsMan.midEval++
			me.eval = EvaluateMidPosition(g)
			return me
		}
	}

	var legalMoves []checkers.Move
	if startDepth == currentDepth { // in starting position - use those moves there optimized with best move first
		legalMoves = startingLegalMoves
	} else {
		legalMoves = g.GetLegalMoves()
	}

	if g.PlrTurn == checkers.BluePlayer {
		me.eval = lowestE

		for _, move := range legalMoves {
			futureGame := g
			(&futureGame).PlayMove(move)
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
			futureGame := g
			(&futureGame).PlayMove(move)
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
