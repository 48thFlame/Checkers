//go:build js && wasm

package main

// air --build.cmd "GOOS=js GOARCH=wasm go build -o static/checkers.wasm"
import (
	"fmt"
	"syscall/js"

	"github.com/48thFlame/Checkers/checkers"
)

func doNotPanicPlease() {
	r := recover()

	if r != nil { // if "panic"ed
		fmt.Println("Recovered from", r)
	}
}

func main() {
	fmt.Println("Hello World!")
	g := checkers.NewGame()
	fmt.Print(g)
	js.Global().Set("newGame", js.FuncOf(newGameWrapper))

	// avoid exiting the program so can call the functions that are exported to JS
	doneCh := make(chan struct{})
	<-doneCh
}

func convertStateToInt(s checkers.GameState) int {
	switch s {
	case checkers.Playing:
		return 0
	case checkers.BlueWon:
		return 1
	case checkers.RedWon:
		return 2
	case checkers.Draw:
		return 3
	default:
		fmt.Printf("ERROR! unexpected checkers.GameState: %#v\n", s)
		return 0
	}
}

func convertBoardToInts(b checkers.Board) []any {
	board := make([]any, 0, 64)

	for _, slot := range b {
		board = append(board, int(slot))
	}

	return board
}

func convertGameToJsObject(g checkers.Game) map[string]any {
	return map[string]any{
		"state":                 convertStateToInt(g.State),
		"plrTurn":               int(g.PlrTurn),
		"board":                 convertBoardToInts(g.Board),
		"turnNumber":            g.TurnNumber,
		"timeSinceExcitingMove": g.TimeSinceExcitingMove,
	}
}

func newGameWrapper(this js.Value, args []js.Value) any {
	defer doNotPanicPlease()

	g := checkers.NewGame()
	return convertGameToJsObject(*g)
}
