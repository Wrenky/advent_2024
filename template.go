package main

import (
	"advent/helpers"
	log "log/slog"
	"strings"
	"time"

	"github.com/samber/lo"
)

// This needs to change to match the input
func parseInput(input string) []string {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) string {
		return line
	})
}

func main() {
	// Handle command line
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)
	parsed := parseInput(data)
	for _, v := range parsed {
		log.Debug("", "line", v)
	}

	// Part 1
	pre1 := time.Now()
	p1 := 0 //work here
	post1 := time.Now()
	log.Info("Part1", "answer", p1, "time", post1.Sub(pre1))

	// Part 2
	pre2 := time.Now()
	p2 := 0 // work here
	post2 := time.Now()
	log.Info("Part2", "answer", p2, "time", post2.Sub(pre2))
}
