package main

import (
	"advent/helpers"
	"advent/helpers/grid"
	"fmt"
	log "log/slog"
	"strings"

	"github.com/samber/lo"
)

func parseInput(input string) grid.Grid[string] {
	return lo.Map(strings.Split(input, "\n"), func(line string, _ int) []string {
		return strings.Split(line, "")
	})
}

var (
	UP    = grid.Coord{X: 0, Y: 1}
	DOWN  = grid.Coord{X: 0, Y: -1}
	LEFT  = grid.Coord{X: -1, Y: 0}
	RIGHT = grid.Coord{X: 1, Y: 0}
)
var Directions = []grid.Coord{UP, DOWN, LEFT, RIGHT}

func findNext(g grid.Grid[string]) (grid.Coord, bool) {
	res := lo.FilterMap(g, func(row []string, x int) (grid.Coord, bool) {
		// Now do the search
		_, y, ok := lo.FindIndexOf(row, func(val string) bool {
			return val != "."
		})
		return grid.Coord{X: x, Y: y}, ok
	})
	if len(res) == 0 {
		return grid.Coord{}, false
	}
	return res[0], true

}

func Perimeter(region []grid.Coord) int {
	return lo.Sum(lo.Map(region, func(p grid.Coord, _ int) int {
		// Number of points inside my region that I touch?
		touching := lo.Sum(lo.Map(Directions, func(dir grid.Coord, _ int) int {
			next := grid.Add(p, dir)
			if lo.Contains(region, next) {
				return 1
			}
			return 0
		}))
		return 4 - touching
	}))

}

func edgeToAdjancy(dir grid.Coord) []grid.Coord {

	switch dir {
	case UP:
		fallthrough
	case DOWN:
		return []grid.Coord{LEFT, RIGHT}
	case LEFT:
		fallthrough
	case RIGHT:
		return []grid.Coord{UP, DOWN}
	default:
		panic("oh no")
	}
}

func makeGroup(curr grid.Coord, dirs, edges, group []grid.Coord) []grid.Coord {
	// modfies edges :/

	a, b := grid.Add(curr, dirs[0]), grid.Add(curr, dirs[1])
	// If a or b are adjacent, group it together and remove from edges.
	ai := lo.IndexOf(edges, a)
	bi := lo.IndexOf(edges, b)

	// Both are bad, we are done
	if ai == -1 && bi == -1 {
		return group
	}
	var l, r []grid.Coord
	if ai != -1 {
		//Real node! remove from edges as it matched, and add to group. Continue.
		newMember := edges[ai]
		edges = append(edges[:ai], edges[ai+1:]...)
		l = makeGroup(newMember, dirs, edges, append(group, newMember))
	}
	if bi != -1 {
		newMember := edges[bi]
		edges = append(edges[:bi], edges[bi+1:]...)
		r = makeGroup(newMember, dirs, edges, append(group, newMember))
	}
	return append(l, r...)
}

func DiscountPerimeter(g grid.Grid[string], region []grid.Coord) int {
	fmt.Printf("Region starting: %v\n", region)
	eGroups := lo.Map(Directions, func(dir grid.Coord, _ int) []grid.Coord {
		fmt.Printf("Edge %s starting\n", dir)
		edges := lo.Filter(region, func(p grid.Coord, _ int) bool {
			next := grid.Add(p, dir)
			fmt.Printf("\tKWREN: (%s) is %t\n", next, grid.InBounds(g, next))
			// If the next in this direction is out of bounds OR does not equal current, then its an edge.
			if grid.InBounds(g, next) {
				fmt.Printf("\tKWREN: in bounds curr %s, next, %s, check %t\n", p, next, (grid.Get(g, p) != grid.Get(g, next)))
			}

			return (!grid.InBounds(g, next)) || (grid.Get(g, p) != grid.Get(g, next))
		})
		fmt.Printf("\tEdges found: %v\n", edges)
		// Have my set of edges, and my direction. Now I need to group them by check adjanceny by direction
		if len(edges) == 0 {
			return []grid.Coord{}
		}
		if len(edges) == 1 {
			return []grid.Coord{edges[0]}
		}
		dirs := edgeToAdjancy(dir)
		//groups on this edge
		res := makeGroup(edges[0], dirs, edges[1:], []grid.Coord{edges[0]})
		fmt.Printf("\tResult of grouping: %v\n", res)
		return res
	})
	for _, ele := range eGroups {
		fmt.Printf("Final group %v\n", ele)
	}
	return 0
}

func FindRegions(startGrid grid.Grid[string]) [][]grid.Coord {
	// Start with a point that isnt "."
	// Check surrounding points, if matching our point, return and continue search
	// at end, mark the region as "discovered"
	// Check all directions (except prev) for curr +1, if 9 return 1

	var f func(grid.Grid[string], [][]grid.Coord) [][]grid.Coord
	f = func(g grid.Grid[string], agg [][]grid.Coord) [][]grid.Coord {
		point, ok := findNext(g)
		if !ok {
			return agg
		}
		region := lo.Uniq(Traverse(point, g, grid.Get(g, point), []grid.Coord{}))

		return f(g, append(agg, region))
	}
	return f(startGrid, [][]grid.Coord{})
}

func Traverse(curr grid.Coord, g grid.Grid[string], region string, agg []grid.Coord) []grid.Coord {

	agg = append(agg, curr)
	g[curr.X][curr.Y] = "."

	valids := lo.FilterMap(Directions, func(dir grid.Coord, _ int) (grid.Coord, bool) {
		next := grid.Add(curr, dir)
		return next, (grid.InBounds(g, next) && (grid.Get(g, next) == region))
	})
	if len(valids) == 0 {
		return agg
	}

	return lo.FlatMap(valids, func(next grid.Coord, _ int) []grid.Coord {
		return Traverse(next, g, region, agg)
	})
}

func main() {
	// Handle command line
	args := helpers.HandleCommandLine()
	data := helpers.ReadFile(args.InputFile)
	parsed := parseInput(data)
	g := parsed

	regions := FindRegions(g)
	cost := lo.Sum(lo.Map(regions, func(r []grid.Coord, _ int) int {
		return len(r) * Perimeter(r)
	}))
	log.Debug("Part1", "cost", cost)
	c2 := lo.Sum(lo.Map(regions, func(r []grid.Coord, _ int) int {
		log.Debug("info", "region", r, "perim", DiscountPerimeter(g, r))
		return len(r) * DiscountPerimeter(g, r)
	}))
	log.Debug("Part2", "cost", c2)
}
