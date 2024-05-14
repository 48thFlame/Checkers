package ai

import (
	"fmt"

	"github.com/48thFlame/Checkers/checkers"
)

func SmartAi(g checkers.Game) checkers.Move {
	if g.State != checkers.Playing {
		// ! should never get here
		return checkers.Move{}
	}

	isBlueMaxingTurn := g.PlrTurn == checkers.BluePlayer

	var bestEval float64

	if isBlueMaxingTurn {
		bestEval = lowestE
	} else {
		bestEval = highestE
	}

	var bestMove checkers.Move

	legalMoves := g.GetLegalMoves()
	for _, move := range legalMoves {
		gameAfterMovePlayed := g
		(&gameAfterMovePlayed).PlayMove(move)
		eval := minMax(gameAfterMovePlayed, 8, lowestE, highestE)

		if isBlueMaxingTurn {
			if eval > bestEval {
				bestEval = eval
				bestMove = move
			}
		} else {
			if eval < bestEval {
				bestEval = eval
				bestMove = move
			}
		}
	}

	fmt.Printf("bestEval: %v\n", bestEval)

	return bestMove
}
