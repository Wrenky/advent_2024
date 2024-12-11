package main

import (
	"advent/helpers"
	"fmt"
	log "log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/samber/lo"
)

// This needs to change to match the input
func parseInput(input string) []int {
	return lo.Map(strings.Fields(input), func(s string, _ int) int {
		return helpers.Atoi(s)
	})
}

type set struct {
	stone int
	count int
}

var cache = make(map[set]int)

func Blink(a, count int) int {
	me := set{a, count}
	// Already solved for this stone at count X
	if v, ok := cache[me]; ok {
		return v
	}

	b := strconv.Itoa(a)
	switch {
	case count == 0:
		// Iteration done, just add one stone
		cache[me] = 1
	case a == 0:
		// Convert it to 1
		cache[me] = Blink(1, count-1)
	case len(b)%2 == 0:
		// Split and call on both ends
		left := helpers.Atoi(b[0 : len(b)/2])
		right := helpers.Atoi(b[len(b)/2:])
		cache[me] = Blink(left, count-1) + Blink(right, count-1)
	default:
		// Multiply by 2024
		cache[me] = Blink(a*2024, count-1)

	}

	return cache[me]
}

func main() {
	// Handle command line
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)
	parsed := parseInput(data)
	run := func(part, count int) {
		pre := time.Now()
		ans := lo.Sum(lo.Map(parsed, func(a int, _ int) int {
			return Blink(a, count)
		}))
		post := time.Now()
		log.Info(fmt.Sprintf("Part%d", part), "answer", ans, "time", post.Sub(pre))
	}
	run(1, 25)
	run(2, 75)
}
