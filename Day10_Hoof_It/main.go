package main

import (
	"advent/helpers"
	"advent/helpers/grid"
	"fmt"
	log "log/slog"
	"strings"
	"time"

	"github.com/samber/lo"
)

// Just to keep things clean
var (
	UP    = grid.Coord{X: 0, Y: 1}
	DOWN  = grid.Coord{X: 0, Y: -1}
	LEFT  = grid.Coord{X: -1, Y: 0}
	RIGHT = grid.Coord{X: 1, Y: 0}
)
var Directions = []grid.Coord{UP, DOWN, LEFT, RIGHT}

// This needs to change to match the input
func parseInput(input string) [][]int {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) []int {
		return lo.Map(strings.Split(line, ""), func(s string, _ int) int {
			return helpers.Atoi(s)
		})
	})
}

func Trailheads(g grid.Grid[int]) []grid.Coord {
	return lo.FlatMap(g, func(r []int, x int) []grid.Coord {
		return lo.FilterMap(r, func(val int, y int) (grid.Coord, bool) {
			return grid.Coord{X: x, Y: y}, val == 0
		})
	})
}

func traverse(g grid.Grid[int], curr, prev grid.Coord, tracker func(grid.Coord) bool) int {

	// Call the tracker to see if we should continue
	if tracker(curr) {
		return 0
	}

	currentVal := grid.Get(g, curr)
	if currentVal == 9 {
		return 1
	}

	// Check all directions (except prev) for curr +1, if 9 return 1
	valids := lo.FilterMap(Directions, func(dir grid.Coord, _ int) (grid.Coord, bool) {
		next := grid.Add(curr, dir)
		return next, (next != prev && grid.InBounds(g, next) && (grid.Get(g, next) == currentVal+1))
	})

	return lo.Sum(lo.Map(valids, func(next grid.Coord, _ int) int {
		return traverse(g, next, curr, tracker)
	}))
}

// Use a closure to keep the map locked, check map on visit
func fetchChecker() func(grid.Coord) bool {
	visited := make(map[grid.Coord]bool)
	return func(curr grid.Coord) bool {
		if _, ok := visited[curr]; ok {
			return true
		}
		visited[curr] = true
		return false
	}
}

// Fill the function, but always return false
func fetchNOP() func(grid.Coord) bool {
	return func(_ grid.Coord) bool {
		return false
	}
}

func main() {
	// Handle command line
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)
	hikingMap := parseInput(data)
	for _, v := range hikingMap {
		log.Debug("", "line", v)
	}
	trailheads := Trailheads(hikingMap)

	// Do the work, timed. Use the fetchers to specialize the run
	run := func(part int, fetcher func() func(grid.Coord) bool) {
		pre := time.Now()
		ans := lo.Sum(lo.Map(trailheads, func(a grid.Coord, _ int) int {
			return traverse(hikingMap, a, a, fetcher())
		}))
		post := time.Now()
		log.Info(fmt.Sprintf("Part%d", part), "answer", ans, "time", post.Sub(pre))
	}
	run(1, fetchChecker)
	run(2, fetchNOP)
}
