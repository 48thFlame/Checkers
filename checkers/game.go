package checkers

type Player uint8

const (
	BluePlayer Player = iota
	RedPlayer
)

type GameState string

const (
	Playing GameState = "Playing..."
	BlueWon GameState = "Yay Blue!"
	RedWon  GameState = "Red!!!!"
	Draw    GameState = "Its a draw."
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

var (
	RedPieces  = [...]BoardSlot{RedPiece, RedKing}
	BluePieces = [...]BoardSlot{BluePiece, BlueKing}
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
}

func NewGame() *Game {
	return &Game{
		State:           Playing,
		PlrTurn:         BluePlayer,
		Board:           NewBoard(),
		NPiecesCaptured: 0,

		TimeSinceExcitingMove: 0,
	}
}

type Game struct {
	State           GameState
	PlrTurn         Player // who's the current player's turn
	Board           Board
	NPiecesCaptured int

	// this is to check draw condition
	TimeSinceExcitingMove int // time since capture/non-king move
}

func (g *Game) PlayMove(m Move) {
	switch g.PlrTurn {
	case BluePlayer:
		if isOnEnd(BluePlayer, m.EndI) { // If just moved to an end - "King Me!"
			g.Board[m.EndI] = BlueKing
		} else {
			g.Board[m.EndI] = g.Board[m.StartI]
		}

		g.PlrTurn = RedPlayer

	case RedPlayer:
		if isOnEnd(RedPlayer, m.EndI) {
			g.Board[m.EndI] = RedKing
		} else {
			g.Board[m.EndI] = g.Board[m.StartI]
		}

		g.PlrTurn = BluePlayer
	}

	if isIn(g.Board[m.StartI], RedPiece, BluePiece) || len(m.CapturedPiecesI) > 0 {
		// if its a piece move or a capture - exciting!
		g.TimeSinceExcitingMove = 0
	} else {
		g.TimeSinceExcitingMove++
	}

	g.Board[m.StartI] = Empty

	for _, i := range m.CapturedPiecesI {
		g.Board[i] = Empty
		g.NPiecesCaptured++
	}

	g.State = g.getGameState()
}

func (g *Game) getGameState() GameState {
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
