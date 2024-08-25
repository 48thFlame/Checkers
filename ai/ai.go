package ai

import (
	"fmt"
	"math/rand"
	"slices"
	"sort"
	"time"

	"github.com/48thFlame/Checkers/checkers"
)

type MoveEval struct {
	Depth int
	Move  checkers.Move
	Eval  int
}

func (me MoveEval) String() string {
	return fmt.Sprintf("(%d| %d,%d |%d)",
		me.Depth, me.Move.StartI, me.Move.EndI, me.Eval)
}

func sameMove(a, b checkers.Move) bool {
	if a.StartI != b.StartI || a.EndI != b.EndI {
		return false
	}

	return slices.Equal(a.CapturedPiecesI, b.CapturedPiecesI)
}

func SmartAiTimeBound(g checkers.Game, timeLimit time.Duration) (me MoveEval) {
	var bestMoveEval MoveEval

	timeLimitCh := time.After(timeLimit)
	stop := make(chan bool, 1) // somehow use the other channel twice?

	go func() {
		legalMoves := g.GetLegalMoves()
		bestMoveEval.Move = legalMoves[0]
		agd := newAiGameData(g)

		for depth := 1; true; depth++ { // keep searching deeper until told to stop
			select {
			case <-stop:
				// told to stop (timeout)
				return

			default:
				bestMoveEval = minMax(agd, legalMoves, depth, depth, lowestE, highestE)

				// after a search re-order `legalMoves` with the new info
				// first delete best move, then insert it at the beginning
				// not sure if this actually does something (I don't think it does)
				// but it's just too depressing to remove this

				legalMoves = slices.DeleteFunc(legalMoves, func(m checkers.Move) bool {
					return sameMove(m, bestMoveEval.Move)
				})

				legalMoves = slices.Insert(legalMoves, 0, bestMoveEval.Move)
			}
		}
	}()

	<-timeLimitCh
	stop <- true

	return bestMoveEval
}

func SmartAiDepthLimited(g checkers.Game, depthLimit int) (me MoveEval) {
	var bestMoveEval MoveEval

	legalMoves := g.GetLegalMoves()
	agd := newAiGameData(g)
	bestMoveEval = minMax(agd, legalMoves, depthLimit, depthLimit, lowestE, highestE)

	return bestMoveEval
}

func CalculateAllMoves(g checkers.Game, depth int) []MoveEval {
	moveEvalsChannel := make(chan MoveEval)

	legalMoves := g.GetLegalMoves()

	for _, move := range legalMoves {
		futureGame := g
		(&futureGame).PlayMove(move)

		go func(m checkers.Move) {
			me := SmartAiDepthLimited(futureGame, depth)
			moveEvalsChannel <- MoveEval{Depth: depth, Move: m, Eval: me.Eval}
		}(move)
	}

	moveEvals := make([]MoveEval, 0)
	for i := 0; i < len(legalMoves); i++ {
		me := <-moveEvalsChannel

		moveEvals = append(moveEvals, me)
	}

	return moveEvals
}

func SortMoveEvalsHighToLow(s []MoveEval) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Eval > s[j].Eval
	})
}

func SortMoveEvalsLowToHigh(s []MoveEval) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].Eval < s[j].Eval
	})
}

type AiDifficultySetting struct {
	DepthLimit                             int
	WorstChance, ThirdChance, SecondChance float32
}

func DifficultySetAi(g checkers.Game, settings AiDifficultySetting) MoveEval {
	moveEvals := CalculateAllMoves(g, settings.DepthLimit)

	if g.PlrTurn == checkers.BluePlayer {
		SortMoveEvalsHighToLow(moveEvals)
	} else {
		SortMoveEvalsLowToHigh(moveEvals)
	}

	nMoves := len(moveEvals)

	moveEvalToPlay := moveEvals[0]

	if nMoves > 2 { // at least 3 option - sometimes shouldn't play best move
		randomNum := rand.Float32()

		if randomNum < settings.WorstChance {
			// play worst move
			moveEvalToPlay = moveEvals[nMoves-1]
		} else if randomNum < settings.ThirdChance {
			moveEvalToPlay = moveEvals[2] // 3rd best move
		} else if randomNum < settings.SecondChance {
			moveEvalToPlay = moveEvals[1] // 2nd move
		}
	}

	return moveEvalToPlay
}
