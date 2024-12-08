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

// This needs to change to match the input
func parseInput(input string) [][]string {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) []string {
		return strings.Split(line, "")
	})
}

func findAttenas(g [][]string) map[string][]grid.Coord {
	attenas := make(map[string][]grid.Coord)
	for ri, row := range g {
		for ci, val := range row {
			if val != "." {
				if v, ok := attenas[val]; ok {
					attenas[val] = append(v, grid.Coord{Y: ri, X: ci})
				} else {
					attenas[val] = append([]grid.Coord{}, grid.Coord{Y: ri, X: ci})
				}
			}
		}
	}
	return attenas
}

func next2points(g [][]string, a, b grid.Coord) []grid.Coord {
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
		return inBounds(g, a)
	})
}

func stepXY(a, b grid.Coord) (int, int) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	gcd := helpers.GCD(int(math.Abs(float64(dx))), int(math.Abs(float64(dy))))
	return (dx / gcd), (dy / gcd)
}

func nextPoints(g [][]string, a, b grid.Coord) []grid.Coord {
	stepX, stepY := stepXY(a, b)

	var f func(grid.Coord, int, int, []grid.Coord) []grid.Coord
	f = func(c grid.Coord, sX, sY int, agg []grid.Coord) []grid.Coord {
		next := grid.Coord{
			X: c.X + sX,
			Y: c.Y + sY,
		}
		if !inBounds(g, next) {
			return agg
		}
		return f(next, sX, sY, append(agg, next))
	}
	bak := f(a, -stepX, -stepY, []grid.Coord{})
	fwd := f(b, stepX, stepY, []grid.Coord{})
	return append(bak, fwd...)
}

// Move to grids?
func inBounds(g [][]string, c grid.Coord) bool {
	if c.X < 0 || c.Y < 0 {
		return false
	}
	if (c.X > len(g[0])-1) || (c.Y > len(g)-1) {
		return false
	}
	return true
}

func ValidNodes(g [][]string, nodes []grid.Coord) []grid.Coord {
	valids := lo.Filter(nodes, func(a grid.Coord, _ int) bool {
		return inBounds(g, a)
	})
	freq := helpers.FrequencyMap[grid.Coord](valids)
	return lo.Keys(freq)
}

func findDiagonals(attenas map[string][]grid.Coord) []lo.Tuple2[grid.Coord, grid.Coord] {
	diags := []lo.Tuple2[grid.Coord, grid.Coord]{}
	for _, v := range attenas {
		if len(v) == 1 {
			continue
		}
		pairs := lo.FlatMap(v, func(a grid.Coord, i int) []lo.Tuple2[grid.Coord, grid.Coord] {
			return lo.Map(v[i+1:], func(b grid.Coord, _ int) lo.Tuple2[grid.Coord, grid.Coord] {
				return lo.T2(a, b)
			})
		})
		diags = append(diags, pairs...)
	}
	return diags
}

func main() {
	// Handle command line
	args := HandleCommandLine()

	// Parse input
	data := helpers.ReadFile(args.InputFile)
	parsed := parseInput(data)
	for _, v := range parsed {
		log.Debug("", "line", v)
	}
	pre1 := time.Now()

	// Setup stuff
	attenas := findAttenas(parsed)
	diags := findDiagonals(attenas)

	nodes := lo.FlatMap(diags, func(ele lo.Tuple2[grid.Coord, grid.Coord], _ int) []grid.Coord {
		a, b := ele.Unpack()
		return next2points(parsed, a, b)
	})
	freq := helpers.FrequencyMap[grid.Coord](nodes)
	valids := lo.Keys(freq)
	p1 := len(valids)
	post1 := time.Now()
	log.Info("Part1", "answer", p1, "time", post1.Sub(pre1))

	// Part 2
	pre2 := time.Now()
	expandedNodes := lo.FlatMap(diags, func(ele lo.Tuple2[grid.Coord, grid.Coord], _ int) []grid.Coord {
		a, b := ele.Unpack()
		return nextPoints(parsed, a, b)
	})
	m2 := lo.OmitBy(attenas, func(k string, v []grid.Coord) bool {
		return len(v) < 1
	})
	attenaNodes := lo.Flatten(lo.Values(m2))
	combined := lo.Uniq(append(expandedNodes, attenaNodes...))
	p2 := len(combined)
	post2 := time.Now()
	log.Info("Part2", "answer", p2, "time", post2.Sub(pre2))
}
