package ai

import (
	"math/rand"

	"github.com/48thFlame/Checkers/checkers"
)

func RandomAi(g checkers.Game) checkers.Move {
	moves := g.GetLegalMoves()
	randomMove := moves[rand.Intn(len(moves))]

	return randomMove
}
