package main

import (
	"math/rand"

	"github.com/48thFlame/Checkers/ai"
	"github.com/48thFlame/Checkers/checkers"
)

// TODO: come up with a better system for AI difficulty
type aiDifficultySetting struct {
	depthLimit                             int
	worstChance, thirdChance, secondChance float32
}

func difficultySetAi(g checkers.Game, settings aiDifficultySetting) ai.MoveEval {
	moveEvals := ai.CalculateAllMoves(g, settings.depthLimit)

	if g.PlrTurn == checkers.BluePlayer {
		ai.SortMoveEvalsHighToLow(moveEvals)
	} else {
		ai.SortMoveEvalsLowToHigh(moveEvals)
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

	return moveEvalToPlay
}

func easyAi(g checkers.Game) ai.MoveEval {
	return difficultySetAi(g,
		aiDifficultySetting{
			depthLimit:  4,
			worstChance: 0.20, thirdChance: 0.28, secondChance: 0.46})
}

func mediumAi(g checkers.Game) ai.MoveEval {
	return difficultySetAi(g,
		aiDifficultySetting{
			depthLimit:  5,
			worstChance: 0.11, thirdChance: 0.23, secondChance: 0.33})
}

func hardAi(g checkers.Game) ai.MoveEval {
	return difficultySetAi(g,
		aiDifficultySetting{
			depthLimit:  7,
			worstChance: 0.06, thirdChance: 0.17, secondChance: 0.25})
}

func extraHardAi(g checkers.Game) ai.MoveEval {
	return difficultySetAi(g,
		aiDifficultySetting{
			depthLimit:  8,
			worstChance: 0.03, thirdChance: 0.1, secondChance: 0.15})
}

func impossibleAi(g checkers.Game) ai.MoveEval {
	return ai.SmartAiDepthLimited(g, 9)
}

func simpleAi(g checkers.Game) ai.MoveEval {
	return ai.SmartAiDepthLimited(g, 5)
}

func getAiDiffFunc(diff int) func(checkers.Game) ai.MoveEval {
	switch diff {
	case 1:
		return easyAi
	case 2:
		return mediumAi
	case 3:
		return hardAi
	case 4:
		return extraHardAi
	case 5:
		return impossibleAi
	case 6:
		return simpleAi
	default:
		return impossibleAi
	}
}
