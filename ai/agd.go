package ai

import "github.com/48thFlame/Checkers/checkers"

func newAiGameData(g checkers.Game) aiGameData {
	agd := aiGameData{g: g}

	for _, spot := range g.Board {
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
	}

	return agd
}

type aiGameData struct {
	g checkers.Game

	nBlue int // number of blue pieces (also kings)
	nRed  int
}

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

func (agd *aiGameData) playMove(move checkers.Move) {
	nOfCaptured := len(move.CapturedPiecesI)

	if nOfCaptured > 0 {
		switch agd.g.PlrTurn {
		case checkers.BluePlayer:
			agd.nRed -= nOfCaptured
		case checkers.RedPlayer:
			agd.nBlue -= nOfCaptured
		}
	}

	agd.g.PlayMove(move)
}

const (
	// piecesCapturedForEndGame = 24 - 8
	piecesLeftFromSideForEndGame = 4
)

func (agd *aiGameData) isInEndGame() bool {
	return agd.nBlue <= piecesLeftFromSideForEndGame ||
		agd.nRed <= piecesLeftFromSideForEndGame
}
