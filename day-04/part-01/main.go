package main

import (
	"os"
	"strings"
)

const Filename = "../input.txt"

// [delta_y,delta_x]
var Directions = [8][2]int{
	{-1, 0},  // N
	{-1, 1},  // NE
	{0, 1},   // E
	{1, 1},   // SE
	{1, 0},   // S
	{1, -1},  // SW
	{0, -1},  // W
	{-1, -1}, // NW
}

func main() {
	bytes, err := os.ReadFile(Filename)
	if err != nil {
		println("there was an issue reading the file", Filename)
		return
	}
	content := string(bytes)
	grid := strings.Split(content, "\n")
	width := len(grid)
	height := len(grid[0])
	liftableBlocks := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if grid[y][x] != '@' {
				continue
			}
			surroundingRolls := 0
			for _, d := range Directions {
				delta_y := d[0]
				delta_x := d[1]
				if y+delta_y < 0 || y+delta_y >= height {
					continue
				}
				if x+delta_x < 0 || x+delta_x >= width {
					continue
				}
				if grid[y+delta_y][x+delta_x] == '@' {
					surroundingRolls++
					if surroundingRolls > 4 {
						break
					}
				}
			}
			if surroundingRolls < 4 {
				liftableBlocks++
			}
		}
	}
	println("Result:", liftableBlocks)
}
