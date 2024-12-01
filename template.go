package main

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/samber/lo"
)

//go:embed demo
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

	for _, v := range parsed {
		fmt.Printf("%v\n", v)
	}
	//pre, ans, post := time.Now(), len(parsed), time.Now()
	//fmt.Printf("Part1 answer: %d, in %s\n", ans, post.Sub(pre))
}

// This needs to change to match your actual input
func parseInput(input string) []string {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) string {
		return line
	})
}
