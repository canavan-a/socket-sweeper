package board

import (
	"fmt"
	"math/rand"
)

type GameBoard struct {
	Grid [][]Voxel
}

func NewGameBoard(x, y, bombs int) (gb GameBoard) {
	board := make([][]Voxel, y)
	for i := 0; i < y; i++ {
		row := make([]Voxel, x)
		board[i] = row
	}
	gb.Grid = board

	gb.SeedBombs(bombs)

	return
}

func (gb *GameBoard) Print() {
	for i := range gb.Grid {
		row := gb.Grid[i]
		for j := range row {
			if row[j].IsOpen {
				if row[j].IsBomb {
					fmt.Print("B")
				} else {
					fmt.Print(row[j].Number)
				}
			} else {
				fmt.Print("#")
			}
		}
		fmt.Println()
	}
}

func (gb *GameBoard) PrintOpen() {
	for i := range gb.Grid {
		row := gb.Grid[i]
		for j := range row {
			if row[j].IsBomb {
				fmt.Print("B")
			} else {
				fmt.Print(row[j].Number)
			}
		}
		fmt.Println()
	}
}

func (gb *GameBoard) Check(x, y int) (valid bool) {
	if x < 0 {
		return
	}
	if y < 0 {
		return
	}
	yMax := len(gb.Grid)
	xMax := len(gb.Grid[0])

	if x >= xMax {
		return
	}

	if y >= yMax {
		return
	}

	return true
}

func (gb *GameBoard) Open(x, y int) {
	gb.Grid[y][x].Open()
}

func (gb *GameBoard) SeedBombs(bombs int) {
	for bombs > 0 {
		y := rand.Intn(len(gb.Grid))
		x := rand.Intn(len(gb.Grid[0]))

		if gb.Grid[y][x].IsBomb {

		} else {
			gb.Grid[y][x].IsBomb = true
			bombs--
		}

	}
}

func (gb *GameBoard) SeedNumbers() {
	for y := range gb.Grid {
		for x := range len(gb.Grid[0]) {

			var count int

			if gb.BombCheck(x+1, y) {
				count++
			}
			if gb.BombCheck(x-1, y) {
				count++
			}
			if gb.BombCheck(x, y+1) {
				count++
			}
			if gb.BombCheck(x, y-1) {
				count++
			}

			if gb.BombCheck(x+1, y+1) {
				count++
			}
			if gb.BombCheck(x-1, y-1) {
				count++
			}
			if gb.BombCheck(x-1, y+1) {
				count++
			}
			if gb.BombCheck(x+1, y-1) {
				count++
			}

			gb.Grid[y][x].Number = count
		}
	}
}

func (gb *GameBoard) BombCheck(x, y int) bool {
	if !gb.Check(x, y) {
		return false
	}

	if gb.Grid[y][x].IsBomb {
		return true
	}

	return false
}
