package main

import (
	"fmt"
	"time"

	"github.com/48thFlame/Checkers/checkers"
)

type gamePlayerFunc func(g checkers.Game) checkers.Move

type timedMoveInputFunc func(g checkers.Game, timeLimit time.Duration) (checkers.Move, string)

func PlaySimpleGame(plr1blue, plr2red gamePlayerFunc) {
	g := checkers.NewGame()
	fmt.Print(g)

	for g.State == checkers.Playing {
		var move checkers.Move

		var plrToGo gamePlayerFunc
		var plrToGoName string

		switch g.PlrTurn {
		case checkers.BluePlayer:
			plrToGo = plr1blue
			plrToGoName = "Blue"
		case checkers.RedPlayer:
			plrToGo = plr2red
			plrToGoName = "Red"
		}

		move = plrToGo(*g)

		fmt.Printf("%s went: %v\n", plrToGoName, move)
		g.PlayMove(move)
		fmt.Print(g)
		fmt.Println("--------------------")
	}
}

func PlayTimedGame(plr1blue, plr2red timedMoveInputFunc, timeLimit time.Duration) {
	g := checkers.NewGame()
	fmt.Print(g)

	for g.State == checkers.Playing {
		var move checkers.Move

		var plrToGo timedMoveInputFunc
		var plrToGoName string

		switch g.PlrTurn {
		case checkers.BluePlayer:
			plrToGo = plr1blue
			plrToGoName = "Blue"
		case checkers.RedPlayer:
			plrToGo = plr2red
			plrToGoName = "Red"
		}

		start := time.Now()

		move, eval := plrToGo(*g, timeLimit)
		elapsed := time.Since(start)

		fmt.Println(eval)

		fmt.Printf("%s took %s\n", plrToGoName, elapsed)

		if elapsed > (timeLimit + 10*time.Millisecond) {
			fmt.Println("!!TIME OUT!!")
			fmt.Printf("%s took more then an extra 10ms!\n", plrToGoName)

			switch g.PlrTurn {
			case checkers.BluePlayer:
				fmt.Println("Red won because of time limit")
				g.State = checkers.RedWon
			case checkers.RedPlayer:
				fmt.Println("Blue won because of time limit")
				g.State = checkers.BlueWon
			}

			return
		}

		fmt.Printf("%s went: %v\n", plrToGoName, move)
		g.PlayMove(move)
		fmt.Print(g)
		fmt.Println("--------------------")
	}
}

func SimulateTimedGame(plr1blue, plr2red timedMoveInputFunc, timeLimit time.Duration) checkers.Game {
	g := checkers.NewGame()

	for g.State == checkers.Playing {
		var move checkers.Move

		var plrToGo timedMoveInputFunc

		switch g.PlrTurn {
		case checkers.BluePlayer:
			plrToGo = plr1blue
		case checkers.RedPlayer:
			plrToGo = plr2red
		}

		start := time.Now()

		move, _ = plrToGo(*g, timeLimit)

		elapsed := time.Since(start)

		if elapsed > (timeLimit + 15*time.Millisecond) {
			fmt.Println("Time Out!")
			fmt.Println(elapsed)
			switch g.PlrTurn {
			case checkers.BluePlayer:
				g.State = checkers.RedWon
			case checkers.RedPlayer:
				g.State = checkers.BlueWon
			}

			return *g
		}

		g.PlayMove(move)
	}

	return *g
}

type TournamentResults struct {
	BlueWins int
	RedWins  int
	Draws    int
}

func (tr TournamentResults) String() string {
	return fmt.Sprintf("BlueWins: %d\nRedWins: %d\nDraws: %d\n",
		tr.BlueWins, tr.RedWins, tr.Draws)
}

// play = not concurrent
func PlayTimedTournament(plr1blue, plr2red timedMoveInputFunc, timeLimit time.Duration, nOfGames int) {
	tr := TournamentResults{}

	for i := 1; i <= nOfGames; i++ {
		fmt.Println("Playing game", i)
		g := SimulateTimedGame(plr1blue, plr2red, timeLimit)
		fmt.Print(g) // print final position

		switch g.State {
		case checkers.BlueWon:
			tr.BlueWins++
		case checkers.RedWon:
			tr.RedWins++
		case checkers.Draw:
			tr.Draws++
		}
	}

	fmt.Print(tr)
}
