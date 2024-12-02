package main

import (
	_ "embed"
	"fmt"
	"math"
	"slices"
	"strings"

	"advent/helpers"

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
	loc1, loc2 := parseInput(data)

	// Part1: Sort the two slices, take the absolute value of the difference and sum
	slices.Sort(loc1)
	slices.Sort(loc2)

	res := lo.Reduce(loc1, func(agg int, one int, i int) int {
		return agg + int(math.Abs(float64(one)-float64(loc2[i])))
	}, 0)

	fmt.Printf("Part1: %d\n", res)

	// Part 2: Generate a frequency map of loc2, then sum the products of each loc1 value by its
	//         frequency in the freq map.
	freq := helpers.FrequencyMap(loc2)
	p2 := lo.Reduce(loc1, func(agg int, val int, _ int) int {
		return agg + (val * freq[val])
	}, 0)
	fmt.Printf("Part2: %d\n", p2)
}

func parseInput(input string) ([]int, []int) {
	vals := lo.Map(strings.Split(input, "\n"), func(line string, _ int) []int {
		vals := lo.Map(strings.Fields(line), func(num string, _ int) int {
			return helpers.Atoi(num)
		})
		return vals
	})
	lefty := lo.Map(vals, func(v []int, _ int) int { return v[0] })
	right := lo.Map(vals, func(v []int, _ int) int { return v[1] })
	if len(lefty) != len(right) {
		panic("Different lenghth lists")
	}
	return lefty, right
}
