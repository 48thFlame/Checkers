package ai

import (
	"fmt"
	"slices"
	"time"

	"github.com/48thFlame/Checkers/checkers"
)

type moveEval struct {
	depth int
	move  checkers.Move
	eval  int
}

func (me moveEval) String() string {
	return fmt.Sprintf("(%d| %d,%d |%d)",
		me.depth, me.move.StartI, me.move.EndI, me.eval)
}

func SmartAi(g checkers.Game) checkers.Move {
	var bestMoveEval moveEval

	timeLimitCh := time.After(time.Millisecond * 200)
	stop := make(chan bool, 1) // somehow use the other channel twice?
	defer func() {
		stop <- true
	}()

	go func() {
		legalMoves := g.GetLegalMoves()
		bestMoveEval.move = legalMoves[0]

		for depth := 1; true; depth++ { // keep searching deeper until told to stop
			select {
			case <-stop:
				// told to stop
				return

			default:
				bestMoveEval = minMax(g, legalMoves, depth, depth, lowestE, highestE)

				// after a search re-order `legalMoves` with the new info
				legalMoves = slices.DeleteFunc(legalMoves, func(m checkers.Move) bool {
					return sameMove(m, bestMoveEval.move)
				})

				legalMoves = slices.Insert(legalMoves, 0, bestMoveEval.move)
			}
		}
	}()

	<-timeLimitCh

	fmt.Println(bestMoveEval)

	return bestMoveEval.move
}
