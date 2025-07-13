package main

import (
	"fmt"
	"time"

	"github.com/48thFlame/Checkers/ai"
	"github.com/48thFlame/Checkers/checkers"
)

//TODO: move terminal playing to its own folder and stuff

func main() {
	g := checkers.NewGame()
	me := ai.SmartAiTimeBound(*g, 100*time.Millisecond)
	fmt.Println(me)

	fmt.Println("-----")
	fmt.Print(ai.MinMaxStatsMan)
}
