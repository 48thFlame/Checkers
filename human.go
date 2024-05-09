package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/48thFlame/Checkers/checkers"
)

func validI(i int) bool {
	return 0 < i && i < checkers.BoardSize
}

func _readHumanMoveInput(g checkers.Game) (checkers.Move, error) {
	// startI, endI := -1, -1
	var startIInput, endIInput string
	fmt.Print("Enter you move (startI endI): ")
	fmt.Scanln(&startIInput, &endIInput)

	startI, err := strconv.Atoi(startIInput)
	if err != nil {
		return checkers.Move{}, err
	}
	endI, err := strconv.Atoi(endIInput)
	if err != nil {
		return checkers.Move{}, err
	}

	if !validI(startI) || !validI(endI) {
		return checkers.Move{}, errors.New("please enter valid coordinates")
	}

	startSlot := g.Board[startI]
	notThatTurn := false
	switch g.PlrTurn {
	case checkers.BluePlayer:
		notThatTurn = !isIn(startSlot, checkers.BluePieces[:]...)
	case checkers.RedPlayer:
		notThatTurn = !isIn(startSlot, checkers.RedPieces[:]...)
	}

	if notThatTurn {
		return checkers.Move{}, errors.New("thats not your piece")
	}

	endSlot := g.Board[endI]

	if endSlot != checkers.Empty {
		return checkers.Move{}, errors.New("sorry you can land there")
	}

	legalMoves := g.GetLegalMoves()
	for _, move := range legalMoves {
		if move.StartI == startI && move.EndI == endI {
			// once finds a match then return
			return move, nil
		}
	}

	// if made it here means no move matches
	return checkers.Move{}, errors.New("sorry no such legal move")
}

func humanMove(g checkers.Game) checkers.Move {
	m, err := _readHumanMoveInput(g)
	if err != nil {
		fmt.Println("Error !!:", err)
		return humanMove(g)
	}

	return m
}
