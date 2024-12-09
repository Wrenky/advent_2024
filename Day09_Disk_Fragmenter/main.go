package main

import (
	"advent/helpers"
	"fmt"
	log "log/slog"
	"slices"
	"strings"
	"time"

	"github.com/samber/lo"
)

// This needs to change to match the input
func parseInput(input string) []int {
	return lo.FlatMap(strings.Split(input, "\n"), func(line string, _ int) []int {
		return lo.Map(strings.Split(line, ""), func(a string, _ int) int {
			return helpers.Atoi(a)
		})
	})
}

func expand(in []int) []int {
	var f func(int, int, int, []int, []int) []int
	f = func(fileLen int, spaceLen int, id int, rest []int, agg []int) []int {
		blk := lo.Map(lo.Range(fileLen), func(_ int, _ int) int { return id })
		spaces := lo.Map(lo.Range(spaceLen), func(_ int, _ int) int { return -1 })
		newAgg := slices.Concat(agg, blk, spaces)
		if len(rest) == 1 {
			// Final file!
			extra := lo.Map(lo.Range(rest[0]), func(_ int, _ int) int { return id + 1 })
			return slices.Concat(newAgg, extra)
		}
		return f(rest[0], rest[1], id+1, rest[2:], newAgg)
	}
	return f(in[0], in[1], 0, in[2:], []int{})
}

func checksum(in []int) int {
	return lo.Reduce(in, func(agg int, a int, i int) int {
		if a == -1 {
			a = 0
		}
		return agg + (a * i)
	}, 0)
}

// Find the last file, returning its size, id and offset
func findLastFile(starting []int, off int) (int, int, bool) {

	in := starting[:off+1]
	id, end, ok := lo.FindLastIndexOf(in, func(a int) bool { return a != -1 })
	if !ok {
		return 0, 0, false
	}

	_, start, ok := lo.FindLastIndexOf(in, func(a int) bool { return a != id })
	if !ok {
		return 0, 0, false
	}
	start++
	return start, (end - start + 1), true
}

func findSpaceFromOffset(in []int, off int) (int, int, bool) {
	start := in[off:]
	_, startSp, ok := lo.FindIndexOf(start, func(a int) bool { return a == -1 })
	if !ok {
		return 0, 0, false
	}

	_, endSp, ok := lo.FindIndexOf(start[startSp:], func(a int) bool { return a != -1 })
	if !ok {
		return 0, 0, false
	}
	return off + startSp, endSp, true
}

func moveRange(in []int, l, u, size int) []int {
	for k := 0; k < size; k++ {
		in[l+k], in[u+k] = in[u+k], in[l+k]
	}
	return in
}

func moveFileToSpace(start []int, off int) ([]int, int) {
	// Find last file
	fi, fSize, ok := findLastFile(start, off)
	if !ok || fSize < 0 {
		return start, off - 1
	}

	// Now find a place to put it!
	var f func([]int, int) []int
	f = func(in []int, spOff int) []int {
		si, spSize, ok := findSpaceFromOffset(in, spOff)
		if !ok || si >= fi || spSize < 0 {
			return in
		}
		if spSize >= fSize {
			// Yay, move our file in
			moveRange(in, si, fi, fSize)
			return in
		}
		return f(in, si+spSize+1)
	}
	// this file is now done?
	return f(start, 0), fi - 1
}

func frag(in []int) []int {
	_, fi, _ := lo.FindIndexOf(in, func(a int) bool { return a == -1 })
	_, li, _ := lo.FindLastIndexOf(in, func(a int) bool { return a != -1 })

	if li <= fi {
		return in[:fi]
	}

	in[fi], in[li] = in[li], in[fi]
	return frag(in)
}

func betterFrag(starting []int) []int {
	// Move a file, then do it again while file offset > 0
	var f func([]int, int) []int
	f = func(in []int, off int) []int {
		newIn, newOff := moveFileToSpace(in, off)
		if newOff <= 0 {
			return newIn
		}
		return f(newIn, newOff)
	}
	return f(starting, len(starting)-1)
}

func main() {
	// Handle command line, parse data
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)
	diskMap := parseInput(data)

	// Expanded map, with -1 as "space" and ids filled into size
	expanded := expand(diskMap)

	// Runner that wraps the timing
	run := func(part int, in []int, fragger func([]int) []int) {
		pre := time.Now()
		p := slices.Clone(in)
		chksum := checksum(fragger(p))
		post := time.Now()
		log.Info(fmt.Sprintf("Part%d", part), "answer", chksum, "time", post.Sub(pre))
	}

	run(1, expanded, frag)
	run(2, expanded, betterFrag)
}
