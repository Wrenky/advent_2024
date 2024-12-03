package main

import (
	"advent/helpers"
	_ "embed"
	"fmt"
	"regexp"
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

	fmt.Printf("Part1 answer: %d\n", exec(data))

	fmt.Printf("Part2 answer: %d\n", execSometimes(data))
}

// This needs to change to match your actual input
func exec(input string) int {
	pattern, _ := regexp.Compile(`mul\((\d+),(\d+)\)`)
	res := pattern.FindAllStringSubmatch(input, -1)
	return lo.Sum(lo.Map(res, func(matches []string, _ int) int {
		return helpers.Atoi(matches[1]) * helpers.Atoi(matches[2])
	}))
}

func execSometimes(input string) int {
	pattern, _ := regexp.Compile(`(mul)\((\d+),(\d+)\)|(do)\(\)|(don't)\(\)`)
	res := pattern.FindAllStringSubmatch(input, -1)
	return ExecHelp(true, res, 0)
}

func ExecHelp(keep bool, matches [][]string, agg int) int {
	if len(matches) == 0 {
		return agg
	}
	match := matches[0]
	if match[1] == "mul" {
		if keep {
			return ExecHelp(keep, matches[1:], agg+(helpers.Atoi(match[2])*helpers.Atoi(match[3])))
		} else {
			return ExecHelp(keep, matches[1:], agg)
		}
	} else if match[4] == "do" {
		return ExecHelp(true, matches[1:], agg)
	} else if match[5] == "don't" {
		return ExecHelp(false, matches[1:], agg)
	}
	return 0
}
