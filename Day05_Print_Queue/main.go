package main

import (
	"advent/helpers"
	_ "embed"
	"fmt"
	"slices"
	"strings"

	"github.com/samber/lo"
)

//go:embed input
var data string

type safety_manual struct {
	pages  [][]int
	before map[int][]int
}

func init() {
	// Strip trailing newline
	data = strings.TrimRight(data, "\n")
	if len(data) == 0 {
		panic("No input file")
	}
}

func main() {
	manual := parseInput(data)

	valid := lo.Filter(manual.pages, func(page []int, _ int) bool {
		return manual.ValidPage(page)
	})
	count := lo.Sum(lo.Map(valid, func(page []int, _ int) int {
		return page[len(page)/2]
	}))
	fmt.Printf("Part1: %d\n", count)

	invalidPages := lo.Filter(manual.pages, func(page []int, _ int) bool {
		return !manual.ValidPage(page)
	})
	part2 := lo.Sum(lo.Map(invalidPages, func(page []int, _ int) int {
		res := manual.SortPage(page)
		return res[len(page)/2]
	}))

	fmt.Printf("Part2: %d\n", part2)

}

func (m safety_manual) ValidPage(page []int) bool {
	result := m.SortPage(page)
	return slices.Equal(page, result)
}

// Sort the page based on our map
func (m safety_manual) SortPage(page []int) []int {
	result := slices.Clone(page)
	slices.SortFunc(result, func(x, y int) int {
		if lo.Contains(m.before[x], y) {
			return -1
		} else {
			return 1
		}
	})
	return result
}

func parseInput(input string) safety_manual {
	chunks := lo.Map(strings.Split(input, "\n\n"), func(chunk string, _ int) string {
		return chunk
	})

	type Rule struct {
		X, Y int
	}
	rules := lo.Map(strings.Split(chunks[0], "\n"), func(line string, _ int) Rule {
		sides := strings.Split(line, "|")
		return Rule{
			X: helpers.Atoi(sides[0]),
			Y: helpers.Atoi(sides[1]),
		}
	})
	before := make(map[int][]int)
	for _, rule := range rules {
		if v, ok := before[rule.X]; ok {
			before[rule.X] = lo.Uniq(append(v, rule.Y))
		} else {
			before[rule.X] = append([]int{}, rule.Y)
		}
	}

	pages := lo.Map(strings.Split(chunks[1], "\n"), func(line string, _ int) []int {
		return lo.Map(strings.Split(line, ","), func(str string, _ int) int {
			return helpers.Atoi(str)
		})
	})

	return safety_manual{
		pages:  pages,
		before: before,
	}
}
