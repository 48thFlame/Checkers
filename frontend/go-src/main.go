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

func getAiMoveWrapper(this js.Value, args []js.Value) any {
	defer doNotPanicPlease()

	game := decodeJsToGame(args[0])

	if game.State == checkers.Playing {
		move := ai.HardAi(game)
		game.PlayMove(move)
	}

	return encodeGameToJs(game)
}
