package ai

// agd = AI Game Data
// this is a type that holds all the AI-game data,
// the AI uses this to calculate forward and make decisions

import (
	"math/rand"

	"github.com/48thFlame/Checkers/checkers"
)

type zHash = uint64

const (
	zNPieceTypes = 4
)

// 4 types of pieces, 64 spots for them (technically only 32 but for convenience)
// zobristKeys is a 2-d array of zobrist keys, where the first index is the piece type
// and the second index is the board slot index
// this is used to calculate the hash of the board position
// this way each piece type has a different zobrist key for each slot its on
type zobristKeys [zNPieceTypes][checkers.BoardSize]zHash

func generateZobristKeys() *zobristKeys {
	zk := &zobristKeys{}

	r := rand.New(rand.NewSource(3_14159265358979323)) // just some seed - any will work just be consistent

	for i := 0; i < zNPieceTypes; i++ {
		for j := 0; j < checkers.BoardSize; j++ {
			zk[i][j] = r.Uint64()
		}
	}

	return zk
}

// ZKI = Zobrist Key Index
func pieceTypeToZKI(piece checkers.BoardSlot) int {
	switch piece {
	case checkers.BluePiece:
		return 0
	case checkers.BlueKing:
		return 1
	case checkers.RedPiece:
		return 2
	case checkers.RedKing:
		return 3
	}

	panic("No Zobrist Key for that piece")
}

// this represents the bounds of a transposition table entry
// it can be exact, lower or upper bounds
// bounds means that the value is at least or at most the value of the entry
// this is used to prune the search tree in the AI
// if the bounds are exact, then the value is the exact value of the entry
type tPosTableBounds uint8

const (
	exactBounds tPosTableBounds = iota
	lowerBounds
	upperBounds
)

type tableEntry struct {
	me     MoveEval
	bounds tPosTableBounds
}

type tPosTableType = map[zHash]tableEntry

func newAiGameData(g checkers.Game) aiGameData {
	agd := aiGameData{
		g:         g,
		zk:        generateZobristKeys(),
		tPosTable: make(tPosTableType),
	}

	for spotI, spot := range g.Board {
		switch spot {
		case checkers.NaS, checkers.Empty:
			continue

		case checkers.BluePiece:
			agd.nBlue++

		case checkers.BlueKing:
			agd.nBlue++

		case checkers.RedPiece:
			agd.nRed++

		case checkers.RedKing:
			agd.nRed++
		}

		agd.updateHash(spot, spotI)
	}

	return agd
}

type aiGameData struct {
	g checkers.Game

	nBlue int // number of blue pieces (also kings)
	nRed  int

	zk        *zobristKeys
	h         zHash         // current position hash
	tPosTable tPosTableType // trans-position table (its a `map`)
}

func (agd *aiGameData) updateHash(piece checkers.BoardSlot, slotI int) {
	agd.h ^= agd.zk[pieceTypeToZKI(piece)][slotI]
}

func (agd *aiGameData) playMove(move checkers.Move) {
	nOfCaptured := len(move.CapturedPiecesI)

	if nOfCaptured > 0 {
		switch agd.g.PlrTurn {
		case checkers.BluePlayer:
			agd.nRed -= nOfCaptured
		case checkers.RedPlayer:
			agd.nBlue -= nOfCaptured
		}

		for _, cI := range move.CapturedPiecesI {
			agd.updateHash(agd.g.Board[cI], cI)
		}
	}

	pieceMoved := agd.g.Board[move.StartI]

	agd.updateHash(pieceMoved, move.StartI) // remove from start
	agd.updateHash(pieceMoved, move.EndI)   // add to end

	agd.g.PlayMove(move)
}

// checks whether canCapture - in which case won't end the search there
func (agd *aiGameData) canCapture() bool {
	if agd.g.State != checkers.Playing {
		return false
	}

	for i, slot := range agd.g.Board {
		if slot == checkers.NaS || slot == checkers.Empty {
			continue
		}

		var directionsToUse []int
		var good bool // is looking maybe at wrong slot because not thats player turn?
		var enemyPieces []checkers.BoardSlot

		switch agd.g.PlrTurn {
		case checkers.BluePlayer:
			directionsToUse, good = checkers.GetDirectionsToUse(checkers.BluePlayer, slot)
			if !good {
				continue
			}

			enemyPieces = checkers.RedPieces[:]
		case checkers.RedPlayer:
			directionsToUse, good = checkers.GetDirectionsToUse(checkers.RedPlayer, slot)
			if !good {
				continue
			}

			enemyPieces = checkers.BluePieces[:]
		}

		captures := checkers.GetCapturesForSlotI(agd.g.Board, i, directionsToUse, enemyPieces)
		if len(captures) > 0 {
			return true // don't care what the captures are, just that canCapture
		}
	}

	return false
}

const (
	piecesLeftFromSideForEndGame = 4
)

func (agd *aiGameData) isInEndGame() bool {
	return agd.nBlue <= piecesLeftFromSideForEndGame ||
		agd.nRed <= piecesLeftFromSideForEndGame
}
