package main

import (
	"advent/helpers/grid"
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

func main() {
	g, start := parseInput(data)
	visited := make(map[grid.Coord]int)
	res := moveGuard(grid.Copy(g), start, up, visited)
	fmt.Printf("Part1: %d\n", len(lo.Keys(res)))

	loops := lo.FilterMap(lo.Keys(res), func(c grid.Coord, _ int) (grid.Coord, bool) {
		if c == start {
			return c, false
		}
		newG := grid.Copy(g)
		newG[c.X][c.Y] = "#"
		visited := make(map[state]int)
		return c, checkLoops(newG, start, up, visited)
	})
	fmt.Printf("Part2: %d\n", len(loops))
}

// Direcitons
var (
	up    = grid.Coord{X: -1, Y: 0}
	down  = grid.Coord{X: 1, Y: 0}
	left  = grid.Coord{X: 0, Y: -1}
	right = grid.Coord{X: 0, Y: 1}
)

type state struct {
	c   grid.Coord
	dir grid.Coord
}

func inBounds(g grid.Grid[string], c grid.Coord) bool {
	if c.X < 0 || c.Y < 0 {
		return false
	}
	if (c.X > len(g[0])-1) || (c.Y > len(g)-1) {
		return false
	}
	return true
}

func Turn90(dir grid.Coord) grid.Coord {
	switch dir {
	case up:
		return right
	case right:
		return down
	case down:
		return left
	case left:
		return up
	default:
		panic("Invalid direction passed in!")
	}
}

func print(g grid.Grid[string], pos grid.Coord) {
	g[pos.X][pos.Y] = "o"
	grid.Print(g)
}

func recordVisit[T comparable](visited map[T]int, pos T) map[T]int {
	if v, ok := visited[pos]; ok {
		visited[pos] = v + 1
	} else {
		visited[pos] = 1
	}
	return visited
}

func checkLoops(g grid.Grid[string], pos grid.Coord, dir grid.Coord, visited map[state]int) bool {
	// now we do the same as moveGuard, but check if our map as seen this state before
	nextPos := grid.Add(pos, dir)
	nextState := state{c: pos, dir: dir}

	if !inBounds(g, nextPos) {
		visited = recordVisit(visited, nextState)
		return false
	}
	if _, ok := visited[nextState]; ok {
		// Loop!
		return true
	}
	if grid.Get(g, nextPos) == "#" {
		newDir := Turn90(dir)
		return checkLoops(g, pos, newDir, visited)
	} else {
		// We are not obstructed, record and move foward
		visited = recordVisit(visited, nextState)
		return checkLoops(g, nextPos, dir, visited)
	}
}

func moveGuard(g grid.Grid[string], pos grid.Coord, dir grid.Coord, visited map[grid.Coord]int) map[grid.Coord]int {
	nextPos := grid.Add(pos, dir)

	// if edge of map, return visited
	if !inBounds(g, nextPos) {
		visited = recordVisit(visited, pos)
		return visited
	}
	if grid.Get(g, nextPos) == "#" {
		newDir := Turn90(dir)
		return moveGuard(g, pos, newDir, visited)
	} else {
		// We are not obstructed, record and move foward
		visited = recordVisit(visited, pos)
		return moveGuard(g, nextPos, dir, visited)
	}
}

// This needs to change to match your actual input
func parseInput(input string) (grid.Grid[string], grid.Coord) {
	g := lo.Map(strings.Split(input, "\n"), func(line string, _ int) []string {
		return strings.Split(line, "")
	})

	starts := lo.FilterMap(g, func(row []string, x int) (grid.Coord, bool) {
		_, y, ok := lo.FindIndexOf(row, func(s string) bool {
			return s == "^"
		})
		return grid.Coord{X: x, Y: y}, ok
	})
	start := starts[0]
	g[start.Y][start.X] = "."

	if len(starts) > 1 {
		panic("more than one start found!")
	}
	return g, starts[0]
}
