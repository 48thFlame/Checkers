package main

import (
	"fmt"
	"math/rand"

	"github.com/48thFlame/Checkers/game"
)

func playUntilCant(g *game.Game) {
	for moves := g.GetLegalMoves(); len(moves) != 0; moves = g.GetLegalMoves() {
		fmt.Println(moves)

		randomMove := moves[rand.Intn(len(moves))]
		fmt.Println(randomMove)
		g.PlayMove(randomMove)
		fmt.Print(g.Board)

	}
	fmt.Println("Done")
	fmt.Println(g.GetLegalMoves())

}

func playNMoves(g *game.Game, n int) {
	var moves []game.Move

	for j := 0; j < n; j++ {
		moves = g.GetLegalMoves()
		fmt.Println(moves)

		randomMove := moves[rand.Intn(len(moves))]
		fmt.Println(randomMove)
		g.PlayMove(randomMove)
		fmt.Print(g.Board)
	}
}

func main() {
	g := game.NewGame()
	fmt.Print(g.Board)

	playUntilCant(g)
	// playNMoves(g, 20)
}
