package main

import (
	"advent/helpers"
	_ "embed"
	"fmt"
	"math"
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
	parsed := parseInput(data)

	p1 := lo.Filter(parsed, func(report []int, _ int) bool {
		return SafetyCheck(report)
	})
	fmt.Printf("Part1: %d\n", len(p1))

	// Now get the bad reports, then run them through the dampner
	bad_reports := lo.Filter(parsed, func(report []int, _ int) bool {
		return !SafetyCheck(report)
	})
	p2 := lo.Filter(bad_reports, func(report []int, _ int) bool {
		return ProblemDampner(report)
	})
	fmt.Printf("Part2: %d\n", len(p1)+len(p2))
}

func ProblemDampner(report []int) bool {
	good := lo.Filter(report, func(_ int, i int) bool {
		return SafetyCheck(helpers.RemoveElement(report, i))
	})
	return len(good) > 0
}

// Go cant pass operators, so... this
func greaterThan(a, b int) bool { return a > b }
func lessThan(a, b int) bool    { return a < b }

func SafetyCheck(report []int) bool {
	if report[0] > report[1] {
		return check(report[0], report[1:], greaterThan)
	}
	return check(report[0], report[1:], lessThan)
}

// Check does the actual work here
type comparer func(int, int) bool

func check(prev int, rep []int, comp comparer) bool {
	if len(rep) == 0 {
		return true
	}
	if comp(prev, rep[0]) {
		diff := int(math.Abs(float64(prev) - float64(rep[0])))
		if (diff >= 1) && (diff <= 3) {
			return check(rep[0], rep[1:], comp)
		}
	}
	return false
}

func parseInput(input string) [][]int {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) []int {
		return lo.Map(strings.Fields(line), func(num string, _ int) int {
			return helpers.Atoi(num)
		})
	})
}
