package main

import (
	"github.com/48thFlame/Checkers/ai"
	"github.com/48thFlame/Checkers/checkers"
)

func easyAi(g checkers.Game) ai.MoveEval {
	return ai.DifficultySetAi(g,
		ai.AiDifficultySetting{
			DepthLimit:  4,
			WorstChance: 0.20, ThirdChance: 0.28, SecondChance: 0.46})
}

func mediumAi(g checkers.Game) ai.MoveEval {
	return ai.DifficultySetAi(g,
		ai.AiDifficultySetting{
			DepthLimit:  5,
			WorstChance: 0.11, ThirdChance: 0.23, SecondChance: 0.33})
}

func hardAi(g checkers.Game) ai.MoveEval {
	return ai.DifficultySetAi(g,
		ai.AiDifficultySetting{
			DepthLimit:  7,
			WorstChance: 0.06, ThirdChance: 0.17, SecondChance: 0.25})
}

func extraHardAi(g checkers.Game) ai.MoveEval {
	return ai.DifficultySetAi(g,
		ai.AiDifficultySetting{
			DepthLimit:  8,
			WorstChance: 0.03, ThirdChance: 0.1, SecondChance: 0.15})
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
