package ai

import (
	"fmt"
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

	timeLimitCh := time.After(time.Millisecond * 500)
	stop := make(chan bool, 1) // TODO: somehow use the other channel twice
	defer func() {
		stop <- true
	}()

	go func() {
		legalMoves := g.GetLegalMoves()

		for depth := 1; true; depth++ { // keep searching deeper until told to stop
			select {
			case <-stop:
				// told to stop
				return

			default:
				bestMoveEval = minMax(g, legalMoves, depth, depth, lowestE, highestE)
				legalMoves = getOrderedLegalMoves(&g, bestMoveEval.move)
			}
		}
	}()

	<-timeLimitCh

	fmt.Println(bestMoveEval)

	return bestMoveEval.move
}
