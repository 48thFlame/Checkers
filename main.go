package main

import (
	"fmt"
	"math/rand"
)

func main() {
	b := NewGame()
	fmt.Print(b.Board)
	fmt.Println(b.GetLegalMoves())
	numOfSteps := 0
	for moves := b.GetLegalMoves(); len(moves) != 0; moves = b.GetLegalMoves() {
		numOfSteps++

		randomMove := moves[rand.Intn(len(moves))]
		fmt.Println(randomMove)
		b.PlayMove(randomMove)
		fmt.Print(b.Board)
	}
	fmt.Println("It took", numOfSteps)
}
