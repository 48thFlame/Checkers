package ai

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
	moveEvals := make([]moveEval, 0)

	legalMoves := g.GetLegalMoves()

	depth-- // because playing a move and then min-max

	for _, move := range legalMoves {
		futureGame := *g
		(&futureGame).PlayMove(move)

		eval := minMax(futureGame, depth, lowestE, highestE)
		moveEvals = append(moveEvals, moveEval{move: move, eval: eval})
	}

	return moveEvals
}

func SmartAi(g checkers.Game) (bestMove checkers.Move) {
	moveEvals := calculateAllMoves(&g, 7)
	sortMoveEvalsHighToLow(moveEvals)

	if g.PlrTurn == checkers.BluePlayer {
		// take first move
		bestMove = moveEvals[0].move
	} else {
		// red wants lowest so take last
		bestMove = moveEvals[len(moveEvals)-1].move
	}

	fmt.Printf("%v>%v\n", moveEvals, bestMove)

	return bestMove
}
