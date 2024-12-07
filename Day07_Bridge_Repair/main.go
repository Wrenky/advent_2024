package main

import (
	"advent/helpers"
	"fmt"
	log "log/slog"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/samber/lo"
)

// -----------------------------------------------------------------------------------------
// Boilerplate
// -----------------------------------------------------------------------------------------
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

func parseInput(input string) [][]int {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) []int {
		halves := strings.Split(line, ":")
		target := []int{helpers.Atoi(halves[0])}
		numbers := lo.Map(strings.Fields(halves[1]), func(val string, _ int) int {
			return helpers.Atoi(val)
		})
		return append(target, numbers...)
	})
}

//-----------------------------------------------------------------------------------------

// Just check everything
func operations(target int, nums []int, curr int, withConcat bool) bool {
	if curr > target {
		return false
	}
	if len(nums) == 0 {
		return curr == target
	}
	next, rest := nums[0], nums[1:]
	cc := false
	if withConcat {
		cc = operations(target, rest, concat(curr, next), true)
	}
	return operations(target, rest, (curr*next), false) ||
		operations(target, rest, (curr+next), false) || cc
}

func concat(a, b int) int {
	return helpers.Atoi(fmt.Sprintf("%v%v", a, b))
}

func main() {
	// Handle command line
	args := HandleCommandLine()

	// Parse input
	data := helpers.ReadFile(args.InputFile)
	parsed := parseInput(data)
	for _, v := range parsed {
		log.Debug("", "line", v)
		if len(v) == 3 {
			log.Info("Got a shorter one", "line", v)
		}
	}

	// Part 1
	pre1 := time.Now()
	validEquations := lo.Filter(parsed, func(eq []int, _ int) bool {
		target, curr, rest := eq[0], eq[1], eq[2:]
		return operations(target, rest, curr, false)
	})
	sum := lo.Reduce(validEquations, func(agg int, eq []int, _ int) int {
		return agg + eq[0]
	}, 0)

	post1 := time.Now()
	log.Info("Part1", "answer", sum, "time", post1.Sub(pre1))

	//---------------------------------------------------------------

	// Part 2
	pre2 := time.Now()
	validEquations2 := lo.Filter(parsed, func(eq []int, _ int) bool {
		target, curr, rest := eq[0], eq[1], eq[2:]
		return operations(target, rest, curr, true)
	})
	sum2 := lo.Reduce(validEquations2, func(agg int, eq []int, _ int) int {
		return agg + eq[0]
	}, 0)
	post2 := time.Now()
	log.Info("Part2", "answer", sum2, "time", post2.Sub(pre2))
}
