package ai

import (
	"fmt"
	"slices"
	"time"

	"github.com/48thFlame/Checkers/checkers"
)

func sameMove(a, b checkers.Move) bool {
	if a.StartI != b.StartI || a.EndI != b.EndI {
		return false
	}

	return slices.Equal(a.CapturedPiecesI, b.CapturedPiecesI)
}

func SmartAi(g checkers.Game, timeLimit time.Duration, printEval bool) checkers.Move {
	var bestMoveEval moveEval

	timeLimitCh := time.After(timeLimit)
	stop := make(chan bool, 1) // somehow use the other channel twice?
	defer func() {
		stop <- true
	}()

	go func() {
		legalMoves := g.GetLegalMoves()
		bestMoveEval.move = legalMoves[0]
		agd := newAiGameData(g)

		for depth := 1; true; depth++ { // keep searching deeper until told to stop
			select {
			case <-stop:
				// told to stop
				return

			default:
				bestMoveEval = minMax(agd, legalMoves, depth, depth, lowestE, highestE)

				// after a search re-order `legalMoves` with the new info
				// first delete best move, then insert it at the beginning

				legalMoves = slices.DeleteFunc(legalMoves, func(m checkers.Move) bool {
					return sameMove(m, bestMoveEval.move)
				})

				legalMoves = slices.Insert(legalMoves, 0, bestMoveEval.move)
			}
		}
	}()

	<-timeLimitCh

	if printEval {
		fmt.Println(bestMoveEval)
	}

	return bestMoveEval.move
}
