//go:build js && wasm

// air --build.cmd "GOOS=js GOARCH=wasm go build -o static/checkers.wasm"

package main

import (
	"fmt"
	"syscall/js"

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
	// js.Global().Set("newGame", js.FuncOf(newGameWrapper))
	js.Global().Set("getAiMove", js.FuncOf(getAiMoveWrapper))
	js.Global().Set("getLegalMoves", js.FuncOf(getLegalMovesWrapper))
	js.Global().Set("makeMove", js.FuncOf(makeMoveWrapper))

	// avoid exiting the program so can call the functions that are exported to JS
	doneCh := make(chan struct{})
	<-doneCh
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

func getLegalMovesWrapper(this js.Value, args []js.Value) any {
	defer doNotPanicPlease()

	game := decodeJsToGame(args[0])
	legalMoves := game.GetLegalMoves()

	jsMoves := make([]any, 0)
	for _, m := range legalMoves {
		jsMoves = append(jsMoves, encodeMoveToJs(m))
	}

	return jsMoves
}

func makeMoveWrapper(this js.Value, args []js.Value) any {
	defer doNotPanicPlease()

	game := decodeJsToGame(args[0])

	if game.State == checkers.Playing {
		jsMove := decodeJsMoveToSimpleMove(args[1])
		var moveToMake checkers.Move

		var foundMatch bool
		legalMoves := game.GetLegalMoves()
		for _, move := range legalMoves {
			if move.StartI == jsMove.startI && move.EndI == jsMove.endI {
				foundMatch = true
				moveToMake = move
				break
			}
		}

		if foundMatch {
			game.PlayMove(moveToMake)
		}
	}

	return encodeGameToJs(game)
}
