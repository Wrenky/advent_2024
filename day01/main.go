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
	slices.Sort(loc1)
	slices.Sort(loc2)

	// Slices are sorted, same length so just use loc2 index
	res := lo.Reduce(loc1, func(agg int, one int, i int) int {
		two := loc2[i]
		return agg + int(math.Abs(float64(one)-float64(two)))
	}, 0)

	fmt.Printf("Part1: %d\n", res)

	freq := FrequencyMap(loc2)

	p2 := lo.Sum(lo.Map(loc1, func(val int, _ int) int {
		return val * freq[val]
	}))
	fmt.Printf("Part1: %d\n", p2)
}

func FrequencyMap(m1 []int) map[int]int {
	res := make(map[int]int)
	for _, val := range m1 {
		if v, ok := res[val]; ok {
			res[val] = v + 1
		} else {
			res[val] = 1
		}
	}
	return res
}

// This needs to change to match your actual input
func parseInput(input string) ([]int, []int) {
	vals := lo.Map(strings.Split(input, "\n"), func(line string, _ int) []int {
		vals := lo.Map(strings.Fields(line), func(num string, _ int) int {
			return helpers.Atoi(num)
		})
		return vals
	})
	lefty := lo.Map(vals, func(v []int, _ int) int { return v[0] })
	right := lo.Map(vals, func(v []int, _ int) int { return v[1] })
	return lefty, right
}
