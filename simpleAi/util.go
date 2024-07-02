package simpleAi

import "sort"

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
