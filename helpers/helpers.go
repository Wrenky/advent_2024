package helpers

import (
	"advent/helpers/grid"
	"fmt"
	"math"
	"slices"
	"strconv"
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

func FrequencyMap[A comparable](m1 []A) map[A]int {
	res := make(map[A]int)
	for _, val := range m1 {
		if v, ok := res[val]; ok {
			res[val] = v + 1
		} else {
			res[val] = 1
		}
	}
	return res
}

// ---------------------------------------------------------------
// Slice helpers
// ---------------------------------------------------------------

// Remove an element without  mutating
func RemoveElement(slice []int, s int) []int {
	newSlice := slices.Clone(slice)
	return append(newSlice[:s], newSlice[s+1:]...)
}

// ---------------------------------------------------------------

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
func PicksInnerPoints(c []grid.Coord) int {
	return Shoelace(c) - (len(c) / 2) + 1
}

// Shoelace foruma  is for finding the area of a polygon given its vertex coordinates
// References:
// https://artofproblemsolving.com/wiki/index.php/Shoelace_Theorem
// https://en.wikipedia.org/wiki/Shoelace_formula
func Shoelace(c []grid.Coord) int {
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
