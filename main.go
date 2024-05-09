package main

import (
	"fmt"

	"github.com/48thFlame/Checkers/ai"
	"github.com/48thFlame/Checkers/checkers"
)

type moveInputFunc func(checkers.Game) checkers.Move

func playGame(plr1blue, plr2red moveInputFunc, shouldPrint bool) *checkers.Game {
	g := checkers.NewGame()

	for g.State == checkers.Playing {
		var move checkers.Move

		switch g.PlrTurn {
		case checkers.BluePlayer:
			move = plr1blue(*g)
		case checkers.RedPlayer:
			move = plr2red(*g)
		}

		if shouldPrint {
			fmt.Print(g)
			fmt.Println("Move:", move)
		}

		g.PlayMove(move)
	}

	fmt.Print(g)

	return g
}

func main() {
	g := checkers.NewGame()

	for n := 1; g.State != checkers.Draw; n++ {
		g = playGame(ai.RandomAi, ai.RandomAi, false)
		fmt.Println(g.TimeSinceExcitingMove)

		fmt.Println("Game num:", n)
	}
}
