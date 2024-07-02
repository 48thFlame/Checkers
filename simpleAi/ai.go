package simpleAi

import (
	"fmt"

	"github.com/48thFlame/Checkers/checkers"
)

type moveEval struct {
	move checkers.Move
	eval int
}

func (me moveEval) String() string {
	return fmt.Sprintf("(%d,%d|%d)",
		me.move.StartI, me.move.EndI, me.eval)
}

func calculateAllMoves(g *checkers.Game, depth uint) []moveEval {
	moveEvalsChannel := make(chan moveEval)

	legalMoves := g.GetLegalMoves()

	depth-- // because playing a move and then min-max

	for _, move := range legalMoves {
		futureGame := *g
		(&futureGame).PlayMove(move)

		go func(m checkers.Move) {
			eval := minMax(futureGame, depth, lowestE, highestE)
			// eval := minMax(futureGame, depth, lowestE, highestE)
			moveEvalsChannel <- moveEval{move: m, eval: eval}
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

func SimpleAi(g checkers.Game) (bestMove checkers.Move) {
	moveEvals := calculateAllMoves(&g, 8)
	sortMoveEvalsHighToLow(moveEvals)

	if g.PlrTurn == checkers.BluePlayer {
		// take first move
		bestMove = moveEvals[0].move
	} else {
		// red wants lowest so take last
		bestMove = moveEvals[len(moveEvals)-1].move
	}

	// fmt.Println("simpleE:", eval)
	// fmt.Printf("%v>%v\n", moveEvals, bestMove)

	return bestMove
}
