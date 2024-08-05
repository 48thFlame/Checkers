package main

import (
	"fmt"
	"time"

	"github.com/48thFlame/Checkers/ai"
	"github.com/48thFlame/Checkers/simpleAi"
)

func main() {
	PlayTournament(ai.SmartAi, simpleAi.SimpleAi, 200*time.Millisecond, 5)

	fmt.Println("-----")
	fmt.Print(ai.MinMaxStatsMan)
}
