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

func calculateAllMoves(g *checkers.Game, depth int, legalMoves []checkers.Move) []moveEval {
	moveEvalsChannel := make(chan moveEval)

	depth-- // because playing a move and then min-max

	for _, move := range legalMoves {
		futureGame := *g
		(&futureGame).PlayMove(move)

		go func(m checkers.Move) {
			eval := minMax(futureGame, depth, depth, lowestE, highestE)
			moveEvalsChannel <- moveEval{depth: depth + 1, move: m, eval: eval}
			// moveEvals = append(moveEvals, moveEval{move: move, eval: eval})
		}(move)
	}

	moveEvals := make([]moveEval, 0)
	for i := 0; i < len(legalMoves); i++ {
		me := <-moveEvalsChannel

		moveEvals = append(moveEvals, me)
	}

	return moveEvals
}

func SmartAi(g checkers.Game) checkers.Move {
	var bestMoveEval moveEval

	timeLimitCh := time.After(time.Millisecond * 200)
	stop := make(chan bool, 1) // TODO: somehow use the other channel twice
	defer func() {
		stop <- true
	}()

	go func() {
		var moveEvals []moveEval
		moves := g.GetLegalMoves()

		for depth := 1; true; depth++ { // keep searching deeper until told to stop
			select {
			case <-stop:
				return

			default:
				moveEvals = calculateAllMoves(&g, depth, moves)

				if g.PlrTurn == checkers.BluePlayer {
					sortMoveEvalsHighToLow(moveEvals)
				} else {
					sortMoveEvalsLowToHigh(moveEvals)
				}

				bestMoveEval = moveEvals[0]

				moves = getMovesFromMoveEvals(moveEvals) // does this do what I think it does?
			}
		}
	}()

	<-timeLimitCh

	fmt.Println(bestMoveEval)

	return bestMoveEval.move
}
