package main

import (
	"advent/helpers"
	"advent/helpers/grid"
	log "log/slog"
	"math"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/samber/lo"
)

type cli struct {
	Debug     bool   `name:"debug" short:"v"`
	Run       bool   `name:"input" short:"r" description:"Runs the file named \"input\""`
	InputFile string `name:"file" short:"f" default:"demo"`
}

func HandleCommandLine() *cli {
	args := &cli{}
	kong.Parse(args,
		kong.Description("Run code"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			//	Compact: true,
		}),
	)
	if args.Debug {
		log.SetLogLoggerLevel(log.LevelDebug)
	}
	if args.Run && args.InputFile == "demo" {
		args.InputFile = "input"
	}
	return args
}

// ------------------------------------------------------------------

func parseInput(input string) grid.Grid[string] {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) []string {
		return strings.Split(line, "")
	})
}

func findAttennas(g grid.Grid[string]) map[string][]grid.Coord {
	attennas := make(map[string][]grid.Coord)
	for ri, row := range g {
		for ci, val := range row {
			if val != "." {
				v := lo.ValueOr(attennas, val, []grid.Coord{})
				attennas[val] = append(v, grid.Coord{Y: ri, X: ci})
			}
		}
	}
	return attennas
}

// These three functions are the main "work" of the problem
// StepXY returns the X/Y step amounts needed to find the next points using bcos
// next2points returns the point immedately before coord a, and after coord b
// nextPoints returns all the points within the grid along those lines
func stepXY(a, b grid.Coord) (int, int) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	gcd := helpers.GCD(int(math.Abs(float64(dx))), int(math.Abs(float64(dy))))
	return (dx / gcd), (dy / gcd)
}
func next2points(g grid.Grid[string], a, b grid.Coord) []grid.Coord {
	stepX, stepY := stepXY(a, b)
	n1 := grid.Coord{
		X: a.X - stepX,
		Y: a.Y - stepY,
	}
	n2 := grid.Coord{
		X: b.X + stepX,
		Y: b.Y + stepY,
	}
	return lo.Filter([]grid.Coord{n1, n2}, func(a grid.Coord, _ int) bool {
		return grid.InBounds(g, a)
	})
}
func nextPoints(g [][]string, a, b grid.Coord) []grid.Coord {
	stepX, stepY := stepXY(a, b)

	var f func(grid.Coord, int, int, []grid.Coord) []grid.Coord
	f = func(c grid.Coord, sX, sY int, agg []grid.Coord) []grid.Coord {
		next := grid.Coord{
			X: c.X + sX,
			Y: c.Y + sY,
		}
		if !grid.InBounds(g, next) {
			return agg
		}
		return f(next, sX, sY, append(agg, next))
	}
	bak := f(a, -stepX, -stepY, []grid.Coord{})
	fwd := f(b, stepX, stepY, []grid.Coord{})
	return append(bak, fwd...)
}

// Type to clean up all those type signatures
type Diagonal lo.Tuple2[grid.Coord, grid.Coord]

// Map over all the map arrays, building up Diagonals of similar attennas
func findDiagonals(attennas map[string][]grid.Coord) []Diagonal {
	return lo.FlatMap(lo.Values(attennas), func(v []grid.Coord, _ int) []Diagonal {
		if len(v) == 1 {
			// single attenna, toss
			return []Diagonal{}
		}
		return lo.FlatMap(v, func(a grid.Coord, i int) []Diagonal {
			// Exclude ourselves, and nodes we previously found with i+1
			return lo.Map(v[i+1:], func(b grid.Coord, _ int) Diagonal {
				return Diagonal(lo.T2(a, b))
			})
		})
	})
}

// These next two do part1/part2 based on the work we already did with diagonals/attennas
func UniqueAntinodes(cityMap grid.Grid[string], diags []Diagonal) []grid.Coord {
	antinodes := lo.FlatMap(diags, func(ele Diagonal, _ int) []grid.Coord {
		a, b := lo.Tuple2[grid.Coord, grid.Coord](ele).Unpack()
		return next2points(cityMap, a, b)
	})
	freq := helpers.FrequencyMap[grid.Coord](antinodes)
	return lo.Keys(freq)
}

func LineAntinodes(cityMap grid.Grid[string], diags []Diagonal, attennas map[string][]grid.Coord) []grid.Coord {
	antinodes := lo.FlatMap(diags, func(ele Diagonal, _ int) []grid.Coord {
		a, b := lo.Tuple2[grid.Coord, grid.Coord](ele).Unpack()
		return nextPoints(cityMap, a, b)
	})
	m2 := lo.OmitBy(attennas, func(k string, v []grid.Coord) bool {
		return len(v) < 1
	})
	attenaNodes := lo.Flatten(lo.Values(m2))
	return lo.Uniq(append(antinodes, attenaNodes...))
}

func main() {
	// Handle command line
	args := HandleCommandLine()

	// Parse input
	data := helpers.ReadFile(args.InputFile)

	pre := time.Now()
	cityMap := parseInput(data)
	for _, v := range cityMap {
		log.Debug("", "line", v)
	}
	attennas := findAttennas(cityMap)
	diags := findDiagonals(attennas)
	post := time.Now()
	log.Info("Setup timing", "time", post.Sub(pre))

	// Part 1
	pre1 := time.Now()
	antinodes := UniqueAntinodes(cityMap, diags)
	post1 := time.Now()
	log.Info("Part1", "answer", len(antinodes), "time", post1.Sub(pre1))

	// Part 2
	pre2 := time.Now()
	line_antinodes := LineAntinodes(cityMap, diags, attennas)
	post2 := time.Now()

	log.Info("Part2", "answer", len(line_antinodes), "time", post2.Sub(pre2))
}
