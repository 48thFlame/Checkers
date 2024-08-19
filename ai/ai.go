package ai

import (
	"math/rand"
	"slices"
	"sort"
	"time"

	"github.com/48thFlame/Checkers/checkers"
)

func sameMove(a, b checkers.Move) bool {
	if a.StartI != b.StartI || a.EndI != b.EndI {
		return false
	}

	return slices.Equal(a.CapturedPiecesI, b.CapturedPiecesI)
}

func SmartAiTimeBound(g checkers.Game, timeLimit time.Duration) (move checkers.Move, eval string) {
	var bestMoveEval moveEval

	timeLimitCh := time.After(timeLimit)
	stop := make(chan bool, 1) // somehow use the other channel twice?

	go func() {
		legalMoves := g.GetLegalMoves()
		bestMoveEval.move = legalMoves[0]
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
					return sameMove(m, bestMoveEval.move)
				})

				legalMoves = slices.Insert(legalMoves, 0, bestMoveEval.move)
			}
		}
	}()

	<-timeLimitCh
	stop <- true

	return bestMoveEval.move, bestMoveEval.String()
}

func SmartAiDepthLimited(g checkers.Game, depthLimit int) (move checkers.Move, eval string) {
	var bestMoveEval moveEval

	legalMoves := g.GetLegalMoves()
	agd := newAiGameData(g)
	bestMoveEval = minMax(agd, legalMoves, depthLimit, depthLimit, lowestE, highestE)

	return bestMoveEval.move, bestMoveEval.String()
}

func calculateAllMoves(g checkers.Game, depth int) []moveEval {
	moveEvalsChannel := make(chan moveEval)

	legalMoves := g.GetLegalMoves()

	for _, move := range legalMoves {
		futureGame := g
		(&futureGame).PlayMove(move)
		futureAGD := newAiGameData(futureGame)
		futureLegalMoves := futureGame.GetLegalMoves()

		go func(m checkers.Move) {
			me := minMax(futureAGD, futureLegalMoves, depth, depth-1, lowestE, highestE)
			moveEvalsChannel <- moveEval{depth: depth, move: m, eval: me.eval}
		}(move)
	}

	moveEvals := make([]moveEval, 0)
	for i := 0; i < len(legalMoves); i++ {
		me := <-moveEvalsChannel

		moveEvals = append(moveEvals, me)
	}

	return moveEvals
}

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

type aiDifficultySetting struct {
	depth                                  int
	worstChance, thirdChance, secondChance float32
}

func difficultySetAi(g checkers.Game, settings aiDifficultySetting) checkers.Move {
	moveEvals := calculateAllMoves(g, settings.depth)

	if g.PlrTurn == checkers.BluePlayer {
		sortMoveEvalsHighToLow(moveEvals)
	} else {
		sortMoveEvalsLowToHigh(moveEvals)
	}

	nMoves := len(moveEvals)

	moveEvalToPlay := moveEvals[0]

	if nMoves > 2 { // at least 3 option - sometimes shouldn't play best move
		randomNum := rand.Float32()

		if randomNum < settings.worstChance {
			// play worst move
			moveEvalToPlay = moveEvals[nMoves-1]
		} else if randomNum < settings.thirdChance {
			moveEvalToPlay = moveEvals[2] // 3rd best move
		} else if randomNum < settings.secondChance {
			moveEvalToPlay = moveEvals[1] // 2nd move
		}
	}

	return moveEvalToPlay.move
}

func EasyAi(g checkers.Game) checkers.Move {
	return difficultySetAi(g,
		aiDifficultySetting{depth: 3, worstChance: 0.15, thirdChance: 0.3, secondChance: 0.55})
}

func MediumAi(g checkers.Game) checkers.Move {
	return difficultySetAi(g,
		aiDifficultySetting{depth: 6, worstChance: 0.08, thirdChance: 0.26, secondChance: 0.47})
}

func HardAi(g checkers.Game) checkers.Move {
	return difficultySetAi(g,
		aiDifficultySetting{depth: 7, worstChance: 0.05, thirdChance: 0.15, secondChance: 0.33})
}

func ExtraHardAi(g checkers.Game) checkers.Move {
	return difficultySetAi(g,
		aiDifficultySetting{depth: 8, worstChance: 0.03, thirdChance: 0.1, secondChance: 0.23})
}

func ImpossibleAi(g checkers.Game) checkers.Move {
	m, _ := SmartAiTimeBound(g, 400*time.Millisecond)
	return m
}
