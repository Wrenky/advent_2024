package grid

import (
	"fmt"
	"math"

	"github.com/samber/lo"
)

type Grid[T comparable] [][]T

// ---------------------------------------------------------------
// Grid/2d array helpers
// ---------------------------------------------------------------
func Transpose[T comparable](grid Grid[T]) Grid[T] {
	xl := len(grid[0])
	yl := len(grid)
	result := make([][]T, xl)
	for i := range result {
		result[i] = make([]T, yl)
	}
	for i := 0; i < xl; i++ {
		for j := 0; j < yl; j++ {
			result[i][j] = grid[j][i]
		}
	}
	return result
}
func Rotate90[T comparable](grid Grid[T]) Grid[T] {
	transposed := Transpose(grid)
	result := [][]T{}
	for _, v := range transposed {
		result = append(result, lo.Reverse(v))
	}
	return result
}

func RotateN90[T comparable](grid Grid[T]) Grid[T] {
	result := [][]T{}
	for _, v := range grid {
		result = append(result, lo.Reverse(v))
	}
	return Transpose(result)
}

// Helper for get operations to avoid
func Get[T comparable](grid Grid[T], a Coord) T {
	return grid[a.X][a.Y]
}

func Print[T comparable](g Grid[T]) {
	for _, v := range g {
		fmt.Printf("%v\n", v)
	}
}

func Copy[T comparable](g Grid[T]) Grid[T] {
	newG := make(Grid[T], len(g))
	for i := range g {
		newG[i] = make([]T, len(g[i]))
		copy(newG[i], g[i])
	}
	return newG
}

func InBounds[T comparable](g Grid[T], c Coord) bool {
	if c.X < 0 || c.Y < 0 {
		return false
	}
	if (c.X > len(g[0])-1) || (c.Y > len(g)-1) {
		return false
	}
	return true
}

// These are annoying because in math its x,y, but in head is [col (y)][row (x)]
// y is up/down, x is left right
type Coord struct {
	X, Y int
}

func (c Coord) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}
func ManhattanDist(a, b Coord) int {
	distance := math.Abs(float64(a.X-b.X)) + math.Abs(float64(a.Y-b.Y))
	return int(distance)
}

func Add(a, b Coord) Coord {
	return Coord{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}
