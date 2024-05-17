package main

import (
	"fmt"
	"math/rand"

	"github.com/48thFlame/Checkers/ai"
	"github.com/48thFlame/Checkers/checkers"
)

type moveInputFunc func(checkers.Game) checkers.Move

func PlayGame(plr1blue, plr2red moveInputFunc, shouldPrint bool) checkers.Game {
	g := checkers.NewGame()

	for g.State == checkers.Playing {
		var move checkers.Move

		if shouldPrint {
			fmt.Print(g)
		}

		switch g.PlrTurn {
		case checkers.BluePlayer:
			move = plr1blue(*g)
		case checkers.RedPlayer:
			move = plr2red(*g)
		}

		if shouldPrint {
			fmt.Println("Move:", move)
		}

		g.PlayMove(move)

	}

	if shouldPrint {
		fmt.Print(g)
	}

	return *g
}

func RandomAi(g checkers.Game) checkers.Move {
	moves := g.GetLegalMoves()
	randomMove := moves[rand.Intn(len(moves))]

	return randomMove
}

func TestAisMultipleGames(plr1blue, plr2red moveInputFunc, n int) {
	var blueWins, redWins, draws int
	for i := 0; i < n; i++ {
		g := PlayGame(plr1blue, plr2red, false)

		switch g.State {
		case checkers.BlueWon:
			blueWins++
		case checkers.RedWon:
			redWins++
		case checkers.Draw:
			draws++
		case checkers.Playing:
			panic("why you here")
		}

		fmt.Print(g) // print final position
	}
	fmt.Printf("blueWins: %v\n", blueWins)
	fmt.Printf("redWins: %v\n", redWins)
	fmt.Printf("draws: %v\n", draws)
}

func main() {
	// g := checkers.NewGame()
	// fmt.Print(g)
	// fmt.Println(ai.MinMax(*g, 4, -100000, 100000))
	// fmt.Println(ai.EvaluatePosition(*g))

	// for n := 1; g.State != checkers.Draw; n++ {
	// 	g = playGame(ai.RandomAi, ai.RandomAi, false)
	// 	fmt.Println(g.TimeSinceExcitingMove)

	// 	fmt.Println("Game num:", n)
	// }
	// playGame(humanMove, humanMove, true)

	// for n := 1; g.State != checkers.Draw; n++ {
	// 	g = playGame(ai.RandomAi, ai.RandomAi, false)
	// 	fmt.Println(g.TimeSinceExcitingMove)

	// 	fmt.Println("Game num:", n)
	// }
	// playGame(randomAi, ai.SmartAi, true)
	PlayGame(ai.SmartAi, RandomAi, true)
	// playGame(ai.SmartAi, humanMove, true)
	// testAisMultipleGames(ai.SmartAi, ai.SmartAi, 10)
	// playGame(ai.SmartAi, ai.SmartAi, true)

}
