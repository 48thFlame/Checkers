package ai

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/48thFlame/Checkers/checkers"
)

func formatInt(n int) string {
	in := strconv.FormatInt(int64(n), 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

type minMaxStats struct {
	alphaBetaBreak int
	extendedSearch int
	midEval        int
	endEval        int
	gameEval       int
	tableHitsExact int
	tableHitsUpper int
	tableHitsLower int
}

func (mms minMaxStats) String() string {
	s := strings.Builder{}

	s.WriteString(fmt.Sprintf("alphaBetaBreak: %v\n", formatInt(mms.alphaBetaBreak)))
	s.WriteString(fmt.Sprintf("extendedSearch: %v\n", formatInt(mms.extendedSearch)))
	s.WriteString(fmt.Sprintf("midEval: %v\n", formatInt(mms.midEval)))
	s.WriteString(fmt.Sprintf("endEval: %v\n", formatInt(mms.endEval)))
	s.WriteString(fmt.Sprintf("gameEval: %v\n", formatInt(mms.gameEval)))
	s.WriteString(fmt.Sprintf("tableHitsExact: %v\n", formatInt(mms.tableHitsExact)))
	s.WriteString(fmt.Sprintf("tableHitsUpper: %v\n", formatInt(mms.tableHitsUpper)))
	s.WriteString(fmt.Sprintf("tableHitsLower: %v\n", formatInt(mms.tableHitsLower)))

	return s.String()
}

var MinMaxStatsMan = minMaxStats{}

// classic min-max with alpha beta pruning
func minMax(agd aiGameData, legalMoves []checkers.Move, startDepth, currentDepth int, alpha, beta int) (me MoveEval) {
	if agd.g.State != checkers.Playing {
		MinMaxStatsMan.gameEval++
		me.Eval = gameOverEval(agd, startDepth, currentDepth)
		return me
	}

	if stored, ok := agd.tPosTable[agd.h]; ok {
		if stored.me.Depth >= currentDepth {
			if stored.bounds == exactBounds {
				MinMaxStatsMan.tableHitsExact++

				return stored.me

			} else if stored.bounds == lowerBounds {
				MinMaxStatsMan.tableHitsLower++

				if stored.me.Eval >= beta {
					return stored.me
				}

				alpha = max(alpha, stored.me.Eval)
			} else { // upperBound
				MinMaxStatsMan.tableHitsUpper++

				if alpha >= stored.me.Eval {
					return stored.me
				}

				beta = min(beta, stored.me.Eval)
			}
		}
	}

	if currentDepth == 0 {
		if agd.canCapture() {
			// should not end calculation here - in middle of capture sequence
			// extend search for 1 more move
			MinMaxStatsMan.extendedSearch++
			return minMax(agd, legalMoves, startDepth+1, 1, alpha, beta)
		}
		// if not in middle of capture sequence regular eval

		if agd.isInEndGame() {
			MinMaxStatsMan.endEval++
			me.Eval = evaluateEndGamePos(agd)
			return me
		} else {
			MinMaxStatsMan.midEval++
			me.Eval = evaluateMidPosition(agd)
			return me
		}
	}

	if agd.g.PlrTurn == checkers.BluePlayer {
		me.Eval = lowestE
		localAlpha := alpha

		for _, move := range legalMoves {
			futureGame := agd
			(&futureGame).playMove(move)
			futureLegalMoves := futureGame.g.GetLegalMoves()
			currentMoveEval := minMax(futureGame, futureLegalMoves, startDepth, currentDepth-1, localAlpha, beta)

			if currentMoveEval.Eval > me.Eval { // if current move is better then previously checked
				me.Move = move // its new best move
				me.Eval = currentMoveEval.Eval
			}

			localAlpha = max(localAlpha, currentMoveEval.Eval)
			if beta <= localAlpha {
				MinMaxStatsMan.alphaBetaBreak++
				break
			}
		}
	} else {
		// reds minimizing turn
		me.Eval = highestE
		localBeta := beta

		for _, move := range legalMoves {
			futureGame := agd
			(&futureGame).playMove(move)
			futureLegalMoves := futureGame.g.GetLegalMoves()
			currentMoveEval := minMax(futureGame, futureLegalMoves, startDepth, currentDepth-1, alpha, localBeta)

			if currentMoveEval.Eval < me.Eval {
				me.Move = move
				me.Eval = currentMoveEval.Eval
			}

			localBeta = min(localBeta, currentMoveEval.Eval)
			if localBeta <= alpha {
				MinMaxStatsMan.alphaBetaBreak++
				break
			}
		}
	}

	var toStoreBounds tPosTableBounds

	if me.Eval <= alpha {
		toStoreBounds = upperBounds
	} else if me.Eval >= beta {
		toStoreBounds = lowerBounds
	} else {
		toStoreBounds = exactBounds
	}

	me.Depth = currentDepth
	agd.tPosTable[agd.h] = tableEntry{me: me, bounds: toStoreBounds}
	return me
}
