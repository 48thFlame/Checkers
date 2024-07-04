package ai

import (
	"fmt"
	"sort"
	"strings"

	"github.com/48thFlame/Checkers/checkers"
	"github.com/fatih/color"
)

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func sortMoveEvalsHighToLow(s []moveEval) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].eval > s[j].eval
	})
}

func boardSlotToString(s checkers.BoardSlot, coord int, value bool) (str string) {
	var prettyC string
	if coord < 10 && coord > 0 {
		prettyC = fmt.Sprintf(" 0%v", coord)
	} else if coord < 0 && coord > -10 {
		prettyC = fmt.Sprintf(" %v", coord)
	} else if coord > 9 {
		prettyC = fmt.Sprintf(" %v", coord)
	} else if coord < -9 {
		prettyC = fmt.Sprintf("%v", coord)
	} else { // coord == 0
		prettyC = fmt.Sprintf(" 0%v", coord)
	}

	switch s {
	case checkers.NaS:
		str = color.New(color.BgWhite, color.Bold).Sprint("    ")

	case checkers.Empty:
		if value {
			str = color.New(color.BgBlack, color.FgRed).Sprintf("%v ", prettyC)

		} else {
			str = color.New(color.BgBlack, color.FgHiBlack).Sprintf("%v ", prettyC)
		}
	}

	return str
}

func PrintHeatMap(hm heatMap, name string) {
	s := strings.Builder{}

	board := checkers.Board{
		checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty,
		checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS,
		checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty,
		checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS,
		checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty,
		checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.
			NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.
			Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS, checkers.Empty, checkers.NaS,
	}

	s.WriteString(name)
	s.WriteString(":\n")

	for rowI := 0; rowI < checkers.BoardSideSize; rowI++ {
		for j := 0; j < 2; j++ {
			for colI := 0; colI < checkers.BoardSideSize; colI++ {
				//  i = cols_num * (row) + col
				i := checkers.BoardSideSize*rowI + colI
				slot := board[i]

				if j == 0 { // if the upper more row of row
					s.WriteString(boardSlotToString(slot, i, false))
				} else {
					s.WriteString(boardSlotToString(slot, hm[i], true))
				}
			}
			s.WriteRune('\n')
		}
	}

	fmt.Print(s.String())
}
