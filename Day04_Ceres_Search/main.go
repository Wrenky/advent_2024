package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/samber/lo"
)

//go:embed input
var data string

func init() {
	// Strip trailing newline
	data = strings.TrimRight(data, "\n")
	if len(data) == 0 {
		panic("No input file")
	}
}

var next = map[string]string{
	"X": "M",
	"M": "A",
	"A": "S",
	"S": "",
}

func main() {
	words := parseInput(data)

	part1 := lo.Sum(lo.Map(words, func(row []string, y int) int {
		return lo.Sum(lo.Map(row, func(col string, x int) int {
			if col == "X" {
				return words.check(x, y)
			}
			return 0
		}))
	}))
	fmt.Printf("Part 1: %d\n", part1)

	part2 := lo.Sum(lo.Map(words, func(row []string, y int) int {
		return lo.Sum(lo.Map(row, func(col string, x int) int {
			if col == "A" {
				return words.check2(x, y)
			}
			return 0
		}))
	}))
	fmt.Printf("Part 2: %d\n", part2)
}

type WordGrid [][]string

// For part 2, we just need to check if the diagonals are MAS
func (words WordGrid) check2(col, row int) int {
	// Now we need to check the diagonals!
	// Top Left -> bottom right

	if !words.valid(row-1, col-1) || !words.valid(row+1, col+1) {
		return 0
	}
	if !words.valid(row-1, col+1) || !words.valid(row+1, col-1) {
		return 0
	}

	tl, br := words[row-1][col-1], words[row+1][col+1]
	if (tl == "M" && br == "S") || (br == "M" && tl == "S") {
		// Yay, continue!
	} else {
		return 0
	}

	tr, bl := words[row-1][col+1], words[row+1][col-1]
	if (tr == "M" && bl == "S") || (tr == "S" && bl == "M") {
		// Yay!
		return 1
	}
	return 0
}

// PART 1
// ------------------------------------------------------
// we have an x!
func (words WordGrid) check(col, row int) int {
	//  This counts all matches orginating from x,y in all directions
	res := len(lo.Filter(lo.RangeFrom(1, 9), func(dir int, _ int) bool {
		if dir == 5 {
			return false
		}
		x1, y1 := nextCoords(row, col, dir)
		return words.checkDirectional(x1, y1, dir, next["X"])
	}))

	return res
}

func (words WordGrid) valid(row, col int) bool {
	if row < 0 || col < 0 {
		return false
	}
	if (row > len(words[0])-1) || (col > len(words)-1) {
		return false
	}
	return true
}

func (words WordGrid) checkDirectional(row, col int, dir int, nx string) bool {
	if !words.valid(col, row) {
		return false
	}
	current := words[row][col]
	if current != nx {
		return false
	}
	if next[current] == "" {
		return true
	}
	row1, col1 := nextCoords(row, col, dir)
	return words.checkDirectional(row1, col1, dir, next[current])
}

// Direction is just an integer to tell us which way we are checking
func nextCoords(row, col, direction int) (int, int) {
	switch direction {
	case 1:
		return row + 1, col - 1
	case 2:
		return row + 1, col
	case 3:
		return row + 1, col + 1
	case 4:
		return row, col - 1
	case 5:
		// dumb?
		return row, col
	case 6:
		return row, col + 1
	case 7:
		return row - 1, col - 1
	case 8:
		return row - 1, col
	case 9:
		return row - 1, col + 1
	}
	return 0, 0
}

func parseInput(input string) WordGrid {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) []string {
		return strings.Split(line, "")
	})
}
