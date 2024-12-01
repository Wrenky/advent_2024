package helpers

import (
	"fmt"
	"math"
	"strconv"

	"github.com/samber/lo"
)

//Graphs: https://github.com/dominikbraun/graph

// Math helpers!
func GCD(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
func LCM(a, b int) int {
	return ((a * b) / GCD(a, b))
}

// Atoi in AOC is usually only used in parsing, and after a regexp/split so you know its an int.
func Atoi(in string) int {
	i, err := strconv.Atoi(in)
	if err != nil {
		panic(fmt.Sprintf("helpers.Atoi recieved non integer string: %s", err))
	}
	return i
}

// ---------------------------------------------------------------
// Grid/2d array helpers
// ---------------------------------------------------------------
func Transpose[S any](slice [][]S) [][]S {
	xl := len(slice[0])
	yl := len(slice)
	result := make([][]S, xl)
	for i := range result {
		result[i] = make([]S, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = slice[j][i]
		}
	}
	return result
}
func Rotate90[S any](slice [][]S) [][]S {
	transposed := Transpose(slice)
	result := [][]S{}
	for _, v := range transposed {
		result = append(result, lo.Reverse(v))
	}
	return result
}
func RotateN90[S any](slice [][]S) [][]S {
	result := [][]S{}
	for _, v := range slice {
		result = append(result, lo.Reverse(v))
	}
	return Transpose(result)
}

// ---------------------------------------------------------------

// coordinates! Mostly for grid problems
// These are annoying because in math its x,y, but in code is [col][row]
type Coord struct {
	X, Y int
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}
func (a Coord) ManhattanDist(b Coord) int {
	distance := math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y))
	return int(distance)
}

func (c Coord) Add(o Coord) Coord {
	return Coord{
		X: c.X + o.X,
		Y: c.Y + o.Y,
	}
}

// These were used in advent day10 part 2 2023
// --------------------------------------------------------------------------------
// Pick's Theorem finds  the area of a polygon based on the inner lattice points and
// the boundry points.
// With shoelace formula you can calculate inner points!
// https://artofproblemsolving.com/wiki/index.php/Pick%27s_Theorem
// https://en.wikipedia.org/wiki/Pick%27s_theorem
func Picks(inner int, border int) int {
	return inner + (border / 2) - 1
}
func PicksInnerPoints(c []Coord) int {
	return Shoelace(c) - (len(c) / 2) + 1
}

// Shoelace foruma  is for finding the area of a polygon given its vertex coordinates
// References:
// https://artofproblemsolving.com/wiki/index.php/Shoelace_Theorem
// https://en.wikipedia.org/wiki/Shoelace_formula
func Shoelace(c []Coord) int {
	sum := 0
	p0 := c[len(c)-1]
	for _, p1 := range c {
		sum += p0.Y*p1.X - p0.X*p1.Y
		p0 = p1
	}
	res := math.Abs(float64(sum / 2))
	return int(res)
}

// --------------------------------------------------------------------------------
