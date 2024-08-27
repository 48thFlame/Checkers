//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/48thFlame/Checkers/checkers"
)

// returns a js acceptable type for `checkers.Game`
// use return value in js function `js.FuncOf`
func encodeGameToJs(g checkers.Game) map[string]any {
	return map[string]any{
		"state":                 string(g.State),
		"plrTurn":               int(g.PlrTurn),
		"board":                 convertBoardToInts(g.Board),
		"turnNumber":            g.TurnNumber,
		"timeSinceExcitingMove": g.TimeSinceExcitingMove,
	}
}

// accepts a js.Value which should be the js object repressing a checkers.Game and well decodes it..
func decodeJsToGame(gObj js.Value) checkers.Game {
	return checkers.Game{
		State:                 checkers.GameState(gObj.Get("state").String()),
		PlrTurn:               checkers.Player(gObj.Get("plrTurn").Int()),
		Board:                 convertIntsToBoard(gObj),
		TurnNumber:            gObj.Get("turnNumber").Int(),
		TimeSinceExcitingMove: gObj.Get("timeSinceExcitingMove").Int(),
	}
}

func convertBoardToInts(b checkers.Board) []any {
	board := make([]any, 0, 64)

	for _, slot := range b {
		board = append(board, int(slot))
	}

	return board
}

func convertIntsToBoard(gObj js.Value) checkers.Board {
	var board checkers.Board

	for i := 0; i < 64; i++ {
		board[i] = checkers.BoardSlot(gObj.Get("board").Index(i).Int())
	}

	return board
}

func encodeMoveToJs(move checkers.Move) map[string]any {
	return map[string]any{
		"endI":   move.EndI,
		"startI": move.StartI,
	}
}
