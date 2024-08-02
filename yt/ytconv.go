package yt

import (
	"github.com/48thFlame/Checkers/checkers"
	"github.com/ytaragin/checkers/pkg/board"
	"github.com/ytaragin/checkers/pkg/game"
	"github.com/ytaragin/checkers/pkg/players"
)

func YTAI1(g checkers.Game) checkers.Move {
	b := board.NewEmptyBoard()

	BoardSideSize := 8

	for rowI := 0; rowI < BoardSideSize; rowI++ {
		for colI := 0; colI < BoardSideSize; colI++ {
			//  i = cols_num * (row) + col
			i := BoardSideSize*rowI + colI
			slot := g.Board[i]

			piece := getYTPiece(slot)
			if piece != nil {
				pos := board.NewPosition(rowI, colI)
				b.SetPiece(pos, piece)
			}
		}
	}

	nextColor := getYTColor(g.PlrTurn)

	ytGame := game.InitGameFromBoard(b, nextColor, g.TimeSinceExcitingMove)

	ytGame.Dump()

	player := players.MCSTPlayer{Color: nextColor}
	ytMove := player.GetMove(ytGame)
	ytCaptured := ytMove.GetJumpedPositions()

	captures := make([]int, len(ytCaptured))
	for i, v := range ytCaptured {
		captures[i] = ytPosToI(v)
	}

	return checkers.Move{
		StartI:          ytPosToI(ytMove.GetStart()),
		EndI:            ytPosToI(ytMove.GetEnd()),
		CapturedPiecesI: captures,
	}
}

func getYTPiece(slot checkers.BoardSlot) *board.Piece {
	switch slot {
	case checkers.BluePiece:
		return board.RedNormalPiece
	case checkers.BlueKing:
		return board.RedKingPiece
	case checkers.RedPiece:
		return board.BlueNormalPiece
	case checkers.RedKing:
		return board.BlueKingPiece
	}

	return nil

}

func ytPosToI(pos *board.Position) int {
	return (checkers.BoardSideSize * pos.Row) + pos.Col

}

func getYTColor(gameColor checkers.Player) board.PieceColor {
	if gameColor == checkers.BluePlayer {
		return board.Red
	}
	return board.Blue
}
