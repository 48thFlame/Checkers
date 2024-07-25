package ai

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/48thFlame/Checkers/checkers"
	"github.com/fatih/color"
)

func sortMoveEvalsHighToLow(s []moveEval) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].eval > s[j].eval
	})
}

func sortMoveEvalsLowToHigh(s []moveEval) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].eval < s[j].eval
	})
}

func getMovesFromMoveEvals(moveEvals []moveEval) []checkers.Move {
	moves := make([]checkers.Move, 0, len(moveEvals))

	for _, me := range moveEvals {
		moves = append(moves, me.move)
	}

	return moves
}

func iAbs(a int) int {
	if a < 0 {
		return -a
	} else {
		return a
	}
}

func getManhattanDist(a, b int) int {
	aCol := a % checkers.BoardSideSize
	aRow := a / checkers.BoardSideSize

	bCol := b % checkers.BoardSideSize
	bRow := b / checkers.BoardSideSize

	deltaCol := iAbs(aCol - bCol)
	deltaRow := iAbs(aRow - bRow)

	return deltaCol + deltaRow
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

func formatInt(n int) string {
	in := strconv.FormatInt(int64(n), 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}

type minMaxStats struct {
	calls          int
	alphaBetaBreak int
	extendedSearch int
	midEval        int
	endEval        int
}

func (mms minMaxStats) String() string {
	s := strings.Builder{}
	s.WriteString(fmt.Sprintf("calls: %v\n", formatInt(mms.calls)))
	s.WriteString(fmt.Sprintf("alphaBetaBreak: %v\n", formatInt(mms.alphaBetaBreak)))
	s.WriteString(fmt.Sprintf("extendedSearch: %v\n", formatInt(mms.extendedSearch)))
	s.WriteString(fmt.Sprintf("midEval: %v\n", formatInt(mms.midEval)))
	s.WriteString(fmt.Sprintf("endEval: %v\n", formatInt(mms.endEval)))

	return s.String()
}

var MinMaxStatsMan = minMaxStats{}
