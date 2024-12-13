package main

import (
	"advent/helpers"
	"advent/helpers/grid"
	"fmt"
	log "log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/samber/lo"
)

type Game struct {
	A grid.Coord
	B grid.Coord
	P grid.Coord
}

// This needs to change to match the input
func parseInput(input string) []Game {
	return lo.Map(strings.Split(input, "\n\n"), func(chunk string, _ int) Game {
		//log.Debug("Chunk", "c", chunk)
		apat, _ := regexp.Compile(`Button\sA:\sX\+(\d+),\sY\+(\d+)`)
		bpat, _ := regexp.Compile(`Button\sB:\sX\+(\d+),\sY\+(\d+)`)
		ppat, _ := regexp.Compile(`Prize:\sX=(\d+),\sY=(\d+)`)
		a := apat.FindAllStringSubmatch(chunk, -1)
		b := bpat.FindAllStringSubmatch(chunk, -1)
		p := ppat.FindAllStringSubmatch(chunk, -1)
		return Game{
			A: grid.Coord{X: helpers.Atoi(a[0][1]), Y: helpers.Atoi(a[0][2])},
			B: grid.Coord{X: helpers.Atoi(b[0][1]), Y: helpers.Atoi(b[0][2])},
			P: grid.Coord{X: helpers.Atoi(p[0][1]), Y: helpers.Atoi(p[0][2])},
		}
	})
}

func SecretInt(f float64) bool {
	iVal := int64(f)
	fVal := float64(iVal)
	return f == fVal
}

// Coordinates make this nasty lol
// Essentially find the common (Ax*By) - (Ay*Bx) and divide that into a shifted Px Py to find the "best" matching terms.
// Previous AOC trick lol
func PressButtons(g Game) int {
	common := float64((g.A.X * g.B.Y) - (g.A.Y * g.B.X))
	a := float64((g.P.X*g.B.Y)-(g.P.Y*g.B.X)) / common
	b := float64((g.P.Y*g.A.X)-(g.P.X*g.A.Y)) / common
	if SecretInt(a) && SecretInt(b) {
		return int(3*a + b)
	}
	return 0
}

func main() {
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)
	games := parseInput(data)

	run := func(part int, offset int) {
		pre := time.Now()
		tokens := lo.Sum(lo.Map(games, func(g Game, _ int) int {
			off := grid.Coord{X: offset, Y: offset}
			g.P = grid.Add(g.P, off)
			return PressButtons(g)
		}))
		post := time.Now()
		log.Info(fmt.Sprintf("Part%d", part), "answer", tokens, "time", post.Sub(pre))

	}

	run(1, 0)
	run(2, 10000000000000)
}
