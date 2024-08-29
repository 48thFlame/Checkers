package main

import (
	"github.com/48thFlame/Checkers/ai"
	"github.com/48thFlame/Checkers/checkers"
)

func easyAi(g checkers.Game) ai.MoveEval {
	return ai.DifficultySetAi(g,
		ai.AiDifficultySetting{
			DepthLimit:  4,
			WorstChance: 0.23, ThirdChance: 0.28, SecondChance: 0.46})
}

func mediumAi(g checkers.Game) ai.MoveEval {
	return ai.DifficultySetAi(g,
		ai.AiDifficultySetting{
			DepthLimit:  5,
			WorstChance: 0.12, ThirdChance: 0.29, SecondChance: 0.37})
}

func hardAi(g checkers.Game) ai.MoveEval {
	return ai.DifficultySetAi(g,
		ai.AiDifficultySetting{
			DepthLimit:  7,
			WorstChance: 0.08, ThirdChance: 0.22, SecondChance: 0.33})
}

func extraHardAi(g checkers.Game) ai.MoveEval {
	return ai.DifficultySetAi(g,
		ai.AiDifficultySetting{
			DepthLimit:  8,
			WorstChance: 0.04, ThirdChance: 0.14, SecondChance: 0.23})
}

func impossibleAi(g checkers.Game) ai.MoveEval {
	return ai.SmartAiDepthLimited(g, 9)
}

func simpleAi(g checkers.Game) ai.MoveEval {
	return ai.SmartAiDepthLimited(g, 4)
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
