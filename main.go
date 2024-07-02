package main

import (
	"fmt"
	"math/rand"

	"github.com/48thFlame/Checkers/checkers"
	"github.com/48thFlame/Checkers/simpleAi"
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
	// fmt.Println("SimpleEval:", simpleAi.EvaluatePosition(*g))
	// fmt.Println("SmartEval:", ai.EvaluatePosition(*g))

	PlayGame(simpleAi.SimpleAi, simpleAi.SimpleAi, true)

	// TestAisMultipleGames(simpleAi.SimpleAi, simpleAi.SimpleAi, 10)
	// PlayGame(HumanMove, simpleAi.SimpleAi, true)
	// TestAisMultipleGames(ai.SmartAi, simpleAi.SimpleAi, 100)
	// PlayGame(ai.SmartAi, ai.SmartAi, true)
	// PlayGame(simpleAi.SimpleAi, ai.SmartAi, true)
	// PlayGame(ai.SmartAi, simpleAi.SimpleAi, true)
}
