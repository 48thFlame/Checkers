package main

import (
	"fmt"
	"time"

	"github.com/48thFlame/Checkers/ai"
	"github.com/48thFlame/Checkers/checkers"
)

func main() {
	g := checkers.NewGame()
	me := ai.SmartAiTimeBound(*g, 100*time.Millisecond)
	fmt.Println(me)

	me2 := ai.DifficultySetAi(*g,
		ai.AiDifficultySetting{DepthLimit: 8, WorstChance: 0, ThirdChance: 0.15, SecondChance: 0.33})
	fmt.Println(me2)

	fmt.Println("-----")
	fmt.Print(ai.MinMaxStatsMan)
}
