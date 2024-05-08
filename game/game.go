package game

// TODO: decide what should actually be global

type Player uint8

const (
	BluePlayer Player = iota
	RedPlayer
)

type GameState uint8

const (
	Playing GameState = iota
	BlueWon
	RedWon
	Draw
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

const (
	BoardSideSize = 8
	BoardSize     = BoardSideSize * BoardSideSize
)

/*
Board is a 2-d arrays, that's represented in a 1-d array.
Given that 0 is the top left corner and going to higher index means right/down,
these are true:

i = (cols_num * row) + col
col = mod(i, cols_num)
row = floor(i / cols_num)
*/
type Board [BoardSize]BoardSlot

// NewBoard() returns an initialized board set-up for a checkers game
func NewBoard() Board {
	return Board{
		// regular set-up
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
	// 	// complex position
	// 	NaS, Empty, NaS, RedKing, NaS, Empty, NaS, BluePiece,
	// 	BluePiece, NaS, Empty, NaS, Empty, NaS, BluePiece, NaS,
	// 	NaS, BluePiece, NaS, Empty, NaS, Empty, NaS, BluePiece,
	// 	BluePiece, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
	// 	NaS, Empty, NaS, Empty, NaS, RedPiece, NaS, RedPiece,
	// 	Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
	// 	NaS, Empty, NaS, RedPiece, NaS, RedPiece, NaS, Empty,
	// 	RedPiece, NaS, RedPiece, NaS, Empty, NaS, RedPiece, NaS,
	// }
	// return Board{
	// 	// capture testing board
	// 	NaS, RedKing, NaS, BluePiece, NaS, RedKing, NaS, Empty,
	// 	Empty, NaS, Empty, NaS, Empty, NaS, BlueKing, NaS,
	// 	NaS, Empty, NaS, Empty, NaS, RedPiece, NaS, Empty,
	// 	Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
	// 	NaS, Empty, NaS, RedKing, NaS, Empty, NaS, Empty,
	// 	Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
	// 	NaS, RedPiece, NaS, Empty, NaS, Empty, NaS, Empty,
	// 	RedPiece, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
	// }
	// return Board{
	// 	// thing
	// 	NaS, Empty, NaS, BluePiece, NaS, BluePiece, NaS, Empty,
	// 	Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
	// 	NaS, Empty, NaS, Empty, NaS, Empty, NaS, Empty,
	// 	Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
	// 	NaS, Empty, NaS, Empty, NaS, Empty, NaS, Empty,
	// 	Empty, NaS, Empty, NaS, Empty, NaS, Empty, NaS,
	// 	NaS, Empty, Empty, Empty, NaS, Empty, NaS, Empty,
	// 	RedPiece, RedKing, Empty, NaS, RedPiece, NaS, Empty, NaS,
	// }
}

func NewGame() *Game {
	return &Game{
		State:   Playing,
		PlrTurn: BluePlayer,
		Board:   NewBoard(),

		TimeSinceExcitingMove: 0,
	}
}

type Game struct {
	State   GameState
	PlrTurn Player // who's the current player's turn
	Board   Board

	// these are to check draw condition
	TimeSinceExcitingMove int // time since capture/non-king move
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

	if isIn(g.Board[m.startI], RedPiece, BluePiece) || len(m.capturedPiecesI) > 0 {
		// if its a piece move or a capture - exciting!
		g.TimeSinceExcitingMove = 0
	} else {
		g.TimeSinceExcitingMove++
	}

	g.Board[m.startI] = Empty

	for _, i := range m.capturedPiecesI {
		g.Board[i] = Empty
	}

	g.State = g.GetGameState()
}

func (g Game) GetGameState() GameState {
	legalMoves := g.GetLegalMoves()
	if len(legalMoves) == 0 {
		// if someone cant make a move - game over, other player wins
		switch g.PlrTurn {
		case BluePlayer:
			return RedWon
		case RedPlayer:
			return BlueWon
		}
	}

	if g.TimeSinceExcitingMove >= 80 { // 40 turns
		return Draw
	}

	return Playing
}
