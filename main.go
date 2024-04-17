package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

/*
   i = cols_num * (row ) + col

   -col = cols_num * row - i
   col = - (cols_num * row - i)
   col = -cols_num * row + i

   col = mod(i, cols_num)

   cols_num * row = i - col
   row = (i - col) / cols_num

   row = floor(i / cols_num)

   ~~~~~~~~~

   -cn-1 -cn -cn+1
   -1 0 1
   cn-1 cn cn+1
*/

const (
	BoardSideSize = 8
	BoardSize     = BoardSideSize * BoardSideSize
)

type BoardSlot uint8

const (
	NotSpot BoardSlot = iota
	Empty
	RedPiece
	BluePiece
)

func (s BoardSlot) String() (str string) {
	switch s {
	case NotSpot:
		str = color.New(color.BgWhite).Sprint("    ")
	case Empty:
		str = color.New(color.BgBlack).Sprint("    ")
	case BluePiece:
		str = color.New(color.BgBlack, color.FgHiBlue).Sprint(" @@ ")
	case RedPiece:
		str = color.New(color.BgBlack, color.FgHiRed, color.Bold).Sprint(" ## ")
	}

	return str
}

type Board [BoardSize]BoardSlot

func NewBoard() Board {
	return Board{
		NotSpot, BluePiece, NotSpot, BluePiece, NotSpot, BluePiece, NotSpot, BluePiece,
		BluePiece, NotSpot, BluePiece, NotSpot, BluePiece, NotSpot, BluePiece, NotSpot,
		NotSpot, BluePiece, NotSpot, BluePiece, NotSpot, BluePiece, NotSpot, BluePiece,
		Empty, NotSpot, Empty, NotSpot, Empty, NotSpot, Empty, NotSpot,
		NotSpot, Empty, NotSpot, Empty, NotSpot, Empty, NotSpot, Empty,
		RedPiece, NotSpot, RedPiece, NotSpot, RedPiece, NotSpot, RedPiece, NotSpot,
		NotSpot, RedPiece, NotSpot, RedPiece, NotSpot, RedPiece, NotSpot, RedPiece,
		RedPiece, NotSpot, RedPiece, NotSpot, RedPiece, NotSpot, RedPiece, NotSpot,
	}
}

// String() returns a pretty-print of the board
func (b Board) String() string {
	s := strings.Builder{}

	for rowI := 0; rowI < BoardSideSize; rowI++ {
		for j := 0; j < 2; j++ {
			for colI := 0; colI < BoardSideSize; colI++ {
				//  i = cols_num * (row) + col
				i := BoardSideSize*rowI + colI
				slot := b[i]

				s.WriteString(slot.String())
			}
			s.WriteRune('\n')
		}
	}

	return s.String()
}

func main() {
	b := NewBoard()
	fmt.Print(b)
}
