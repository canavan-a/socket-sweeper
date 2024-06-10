package main

import (
	"fmt"
	"main/board"
)

func main() {
	fmt.Println("socket sweeper started")

	gb := board.NewGameBoard(20, 10, 50)

	gb.Print()

	gb.Open(1, 5)
	fmt.Println("")

	gb.PrintOpen()

	gb.SeedNumbers()

	fmt.Println("")

	gb.PrintOpen()

}
