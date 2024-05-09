package checkers

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// TODO: make clean
// coord - needs to know location to display the number, -1 means should not display the coord
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

// String() returns a pretty-print of the board
func (g Game) String() string {
	s := strings.Builder{}

	for rowI := 0; rowI < BoardSideSize; rowI++ {
		for j := 0; j < 2; j++ {
			for colI := 0; colI < BoardSideSize; colI++ {
				//  i = cols_num * (row) + col
				i := BoardSideSize*rowI + colI
				slot := g.Board[i]

				if j == 0 { // if the upper more row of row
					s.WriteString(slot.String(i))
				} else {
					s.WriteString(slot.String(-1))
				}
			}
			s.WriteRune('\n')
		}
	}

	s.WriteString(string(g.State))
	s.WriteRune('\n')
	return s.String()
}
