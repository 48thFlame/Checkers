package game

// TODO: decide what should actually be global

type Player uint8

const (
	BluePlayer Player = iota
	RedPlayer
)

const (
	BoardSideSize = 8
	BoardSize     = BoardSideSize * BoardSideSize
)

type BoardSlot uint8

const (
	NaS   BoardSlot = iota // Not A Spot (a light square)
	Empty                  // an unoccupied dark square
	BluePiece
	BlueKing
	RedPiece
	RedKing
)

/*
Board is a 2-d arrays, that's represented in a 1-d array.
Given that 0 is the top left corner and going to higher index means right/down,
these are true:

i = cols_num * (row ) + col
col = mod(i, cols_num)
row = floor(i / cols_num)
*/
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
	// 	NaS, Empty, NaS, BluePiece, NaS, BluePiece, NaS, Empty,
	// 	Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
	// 	NaS, Empty, NaS, Empty, NaS, Empty, NaS, Empty,
	// 	Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
	// 	NaS, Empty, NaS, Empty, NaS, Empty, NaS, Empty,
	// 	Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
	// 	NaS, Empty, NaS, Empty, NaS, Empty, NaS, Empty,
	// 	RedPiece, NaS, Empty, NaS, RedPiece, NaS, Empty, NaS,
	// }
}

func NewGame() *Game {
	return &Game{
		Board:   NewBoard(),
		PlrTurn: BluePlayer,
	}
}

type Game struct {
	PlrTurn Player // who's the current player's turn
	Board   Board
}

func (g *Game) PlayMove(m Move) {
	switch g.PlrTurn {
	case BluePlayer:
		if isOnEnd(BluePlayer, m.endI) { // If just moved to an end - "King Me!"
			g.Board[m.endI] = BlueKing
		} else {
			g.Board[m.endI] = g.Board[m.startI]
		}

		g.PlrTurn = RedPlayer

	case RedPlayer:
		if isOnEnd(RedPlayer, m.endI) {
			g.Board[m.endI] = RedKing
		} else {
			g.Board[m.endI] = g.Board[m.startI]
		}

		g.PlrTurn = BluePlayer
	}

	g.Board[m.startI] = Empty

	for _, i := range m.capturedPiecesI {
		g.Board[i] = Empty
	}
}
