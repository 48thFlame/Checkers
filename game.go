package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func NewGame() Game {
	return Game{
		Board: NewBoard(),
		Turn:  BluePlayer,
	}
}

type Game struct {
	Turn  Player
	Board Board
}

type Player uint8

const (
	BluePlayer Player = iota
	RedPlayer
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

type BoardSlot uint8

const (
	NaS   BoardSlot = iota // Not A Spot (a light square)
	Empty                  // an unoccupied dark square
	BluePiece
	BlueKing
	RedPiece
	RedKing
)

// TODO: make clean
func (s BoardSlot) String(coord int) (str string) {
	if coord < 0 {
		switch s {
		case NaS:
			str = color.New(color.BgWhite, color.Bold).Sprint("    ")
		case Empty:
			str = color.New(color.BgBlack).Sprint("    ")
		case BluePiece:
			str = color.New(color.BgBlack, color.FgHiCyan).Sprint(" @@ ")
		case BlueKing:
			str = color.New(color.BgBlack, color.FgHiCyan).Sprint(" kk ")
		case RedPiece:
			str = color.New(color.BgBlack, color.FgHiRed).Sprint(" ## ")
		case RedKing:
			str = color.New(color.BgBlack, color.FgHiRed).Sprint(" KK ")
		}
	} else { // meaning should display coord
		var prettyC string
		if coord < 10 {
			prettyC = fmt.Sprintf("0%v", coord)
		} else {
			prettyC = fmt.Sprint(coord)
		}
		switch s {
		case NaS:
			// str = color.New(color.BgWhite, color.Bold).Sprintf("%v  ", prettyC)
			str = color.New(color.BgWhite, color.Bold).Sprint("    ")

		case Empty:
			str = color.New(color.BgBlack).Sprintf("%v  ", prettyC)
		case BluePiece:
			str = color.New(color.BgBlack, color.FgHiCyan).Sprintf("%v@ ", prettyC)
		case BlueKing:
			str = color.New(color.BgBlack, color.FgHiCyan).Sprintf("%vk ", prettyC)
		case RedPiece:
			str = color.New(color.BgBlack, color.FgHiRed).Sprintf("%v# ", prettyC)
		case RedKing:
			str = color.New(color.BgBlack, color.FgHiRed).Sprintf("%vK ", prettyC)
		}
	}

	return str
}

const (
	BoardSideSize = 8
	BoardSize     = BoardSideSize * BoardSideSize
)

type Board [BoardSize]BoardSlot

// NewBoard() returns an initialized board set-up for a checkers game
func NewBoard() Board {
	return Board{
		NaS, BluePiece, NaS, BluePiece, NaS, BluePiece, NaS, BluePiece,
		BluePiece, NaS, BluePiece, NaS, BluePiece, NaS, BluePiece, NaS,
		NaS, BluePiece, NaS, BluePiece, NaS, BluePiece, NaS, BluePiece,
		Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
		NaS, Empty, NaS, Empty, NaS, Empty, NaS, Empty,
		RedPiece, NaS, RedPiece, NaS, RedPiece, NaS, RedPiece, NaS,
		NaS, RedPiece, NaS, RedPiece, NaS, RedPiece, NaS, RedPiece,
		RedPiece, NaS, RedPiece, NaS, RedPiece, NaS, RedPiece, NaS,
	}
	// return Board{
	// 	NaS, BluePiece, NaS, Empty, NaS, Empty, NaS, Empty,
	// 	Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
	// 	NaS, Empty, NaS, Empty, NaS, Empty, NaS, Empty,
	// 	Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
	// 	NaS, Empty, NaS, Empty, NaS, Empty, NaS, Empty,
	// 	Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
	// 	NaS, Empty, NaS, Empty, NaS, Empty, NaS, Empty,
	// 	Empty, NaS, Empty, NaS, RedPiece, NaS, Empty, NaS,
	// }
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

				if j == 0 { // if the upper more row of row
					s.WriteString(slot.String(i))
				} else {
					s.WriteString(slot.String(-1))
				}
			}
			s.WriteRune('\n')
		}
	}

	return s.String()
}

/*
i=0(1) is top right

going up-right is:

  - 7,
    going up-left is:

  - 9,

    going down-right is:

  - 9,
    going down-left is:

  - 7,

Safe to check also edges because if its really an edge and cant go to a side,
after calculation its ends up in a Nas.
*/
const (
	upLeftCalc  = -9
	upRightCalc = -7

	downLeftCalc  = 7
	downRightCalc = 9
)

type Move struct {
	startI int
	endI   int
	// whether or not the piece who moved is king
	// (not if should become king, but rather was he king in the beginning)
	king bool

	capturedPiecesI []int
}

// GetLegalMoves returns a slice all legal moves in position
func (g Game) GetLegalMoves() []Move {
	moves := make([]Move, 0)

	// can start at 1, because 0 (top right) is always a NaS, same reason stops 1 early
	for i := 1; i < BoardSize-1; i++ {
		slot := g.Board[i]

		if slot == NaS || slot == Empty {
			continue
		}

		if (g.Turn == BluePlayer && slot == BlueKing) ||
			(g.Turn == RedPlayer && slot == RedKing) {

			upLeftI := i + upLeftCalc
			upRightI := i + upRightCalc
			downLeftI := i + downLeftCalc
			downRightI := i + downRightCalc

			if downLeftI < BoardSize && g.Board[downLeftI] == Empty {
				moves = append(moves, Move{startI: i, endI: downLeftI, king: true})
			}

			if downRightI < BoardSize && g.Board[downRightI] == Empty {
				moves = append(moves, Move{startI: i, endI: downRightI, king: true})
			}

			if upLeftI > 0 && g.Board[upLeftI] == Empty {
				moves = append(moves, Move{startI: i, endI: upLeftI, king: true})
			}

			if upRightI > 0 && g.Board[upRightI] == Empty {
				moves = append(moves, Move{startI: i, endI: upRightI, king: true})
			}

		} else if g.Turn == BluePlayer && slot == BluePiece {
			// If it's blue, it's top going **down**
			downLeftI := i + downLeftCalc
			downRightI := i + downRightCalc

			if downLeftI < BoardSize && g.Board[downLeftI] == Empty {
				moves = append(moves, Move{startI: i, endI: downLeftI})
			}

			if downRightI < BoardSize && g.Board[downRightI] == Empty {
				moves = append(moves, Move{startI: i, endI: downRightI})
			}
		} else if g.Turn == RedPlayer && slot == RedPiece {
			// If it's red, it's bottom going **up**
			upLeftI := i + upLeftCalc
			upRightI := i + upRightCalc

			if upLeftI > 0 && g.Board[upLeftI] == Empty {
				moves = append(moves, Move{startI: i, endI: upLeftI})
			}

			if upRightI > 0 && g.Board[upRightI] == Empty {
				moves = append(moves, Move{startI: i, endI: upRightI})
			}
		}
	}

	return moves
}

var (
	blueBoardEnd = [...]int{56, 58, 60, 62}
	redBoardEnd  = [...]int{1, 3, 5, 7}
)

func isOnEnd(plr Player, i int) bool {
	switch plr {
	case BluePlayer:
		if isInSlice(i, blueBoardEnd[:]) {
			return true
		}
	case RedPlayer:
		if isInSlice(i, redBoardEnd[:]) {
			return true
		}
	}

	return false
}

func (g *Game) PlayMove(m Move) {
	g.Board[m.startI] = Empty

	switch g.Turn {
	case BluePlayer:
		if m.king || isOnEnd(BluePlayer, m.endI) { // If just moved to an end - "King Me!"
			g.Board[m.endI] = BlueKing
		} else {
			g.Board[m.endI] = BluePiece
		}

		g.Turn = RedPlayer

	case RedPlayer:
		if m.king || isOnEnd(RedPlayer, m.endI) {
			g.Board[m.endI] = RedKing
		} else {
			g.Board[m.endI] = RedPiece
		}

		g.Turn = BluePlayer
	}

	for _, i := range m.capturedPiecesI {
		g.Board[i] = Empty
	}
}
