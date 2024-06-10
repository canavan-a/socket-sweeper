package main

import (
	"fmt"
	"main/board"
)

func main() {
	fmt.Println("socket sweeper started")

	gb := board.NewGameBoard(20, 10, 10)

	gb.Open(1, 5)

	gb.Open(9, 5)

	gb.Open(13, 8)

	gb.Open(2, 4)

	gb.Open(0, 0)

	gb.Open(19, 9)
	// gb.PrintOpen()

	// fmt.Println("")

	// gb.PrintOpen()

	output := gb.ToOutputValue()
	for _, row := range output {
		fmt.Println(row)
	}
	// fmt.Println(output)
	gb.Print()
	gb.WinCheck()

}
