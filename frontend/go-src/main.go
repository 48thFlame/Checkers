//go:build js && wasm

// air --build.cmd "GOOS=js GOARCH=wasm go build -o static/checkers.wasm"

package main

import (
	"fmt"
	"syscall/js"

	"github.com/48thFlame/Checkers/ai"
	"github.com/48thFlame/Checkers/checkers"
)

// defer this function for every function that should be used with `js.FuncOf`
// so even if the function panics the go runtime will keep running and that way
// you can keep calling functions
func doNotPanicPlease() {
	r := recover()

	if r != nil { // if "panic"ed
		fmt.Println("Recovered from", r)
	}
}

func main() {
	js.Global().Set("newGame", js.FuncOf(newGameWrapper))
	js.Global().Set("getAiMove", js.FuncOf(getAiMoveWrapper))

	// avoid exiting the program so can call the functions that are exported to JS
	doneCh := make(chan struct{})
	<-doneCh
}

func newGameWrapper(this js.Value, args []js.Value) any {
	defer doNotPanicPlease()

	g := checkers.NewGame()
	return encodeGameToJs(*g)
}

func easyAi(g checkers.Game) ai.MoveEval {
	return ai.DifficultySetAi(g,
		ai.AiDifficultySetting{
			DepthLimit:  4,
			WorstChance: 0.23, ThirdChance: 0.28, SecondChance: 0.46})
}

func mediumAi(g checkers.Game) ai.MoveEval {
	return ai.DifficultySetAi(g,
		ai.AiDifficultySetting{
			DepthLimit:  5,
			WorstChance: 0.12, ThirdChance: 0.29, SecondChance: 0.37})
}

func hardAi(g checkers.Game) ai.MoveEval {
	return ai.DifficultySetAi(g,
		ai.AiDifficultySetting{
			DepthLimit:  7,
			WorstChance: 0.08, ThirdChance: 0.22, SecondChance: 0.33})
}

func extraHardAi(g checkers.Game) ai.MoveEval {
	return ai.DifficultySetAi(g,
		ai.AiDifficultySetting{
			DepthLimit:  8,
			WorstChance: 0.04, ThirdChance: 0.14, SecondChance: 0.23})
}

func impossibleAi(g checkers.Game) ai.MoveEval {
	return ai.SmartAiDepthLimited(g, 9)
}

func simpleAi(g checkers.Game) ai.MoveEval {
	return ai.SmartAiDepthLimited(g, 4)
}

func getAiDiffFunc(diff int) func(checkers.Game) ai.MoveEval {
	switch diff {
	case 1:
		return easyAi
	case 2:
		return mediumAi
	case 3:
		return hardAi
	case 4:
		return extraHardAi
	case 5:
		return impossibleAi
	case 6:
		return simpleAi
	default:
		return impossibleAi
	}
}

func getAiMoveWrapper(this js.Value, args []js.Value) any {
	defer doNotPanicPlease()

	game := decodeJsToGame(args[0])

	if game.State == checkers.Playing {
		difficulty := args[1].Int()

		aiFunc := getAiDiffFunc(difficulty)
		me := aiFunc(game)
		fmt.Println(me)
		game.PlayMove(me.Move)
	}

	return encodeGameToJs(game)
}
