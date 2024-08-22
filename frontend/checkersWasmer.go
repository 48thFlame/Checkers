package main

import (
	"fmt"

	"github.com/48thFlame/Checkers/checkers"
)

func main() {
	fmt.Println("Hello World!")
	g := checkers.NewGame()
	fmt.Print(g)
}
