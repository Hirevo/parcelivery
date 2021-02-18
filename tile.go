package main

import (
	"sort"

	astar "github.com/beefsack/go-astar"
	funk "github.com/thoas/go-funk"
)

type (
	TileKind int

	Tile struct {
		env       *Environment
		Kind      TileKind
		Position  Position
		Transport *Transport
		Parcel    *Parcel
	}
)

const (
	TKFree TileKind = iota
	TKTransport
	TKParcel
)

func (self Tile) Free() bool {
	return self.Kind == TKFree
}

func (tile Tile) PathNeighbors() []astar.Pather {
	var available []astar.Pather
	if target := tile.env.At(Position{X: tile.Position.X - 1, Y: tile.Position.Y}); target != nil && target.Free() {
		available = append(available, *target)
	}
	if target := tile.env.At(Position{X: tile.Position.X + 1, Y: tile.Position.Y}); target != nil && target.Free() {
		available = append(available, *target)
	}
	if target := tile.env.At(Position{X: tile.Position.X, Y: tile.Position.Y - 1}); target != nil && target.Free() {
		available = append(available, *target)
	}
	if target := tile.env.At(Position{X: tile.Position.X, Y: tile.Position.Y + 1}); target != nil && target.Free() {
		available = append(available, *target)
	}
	return available
}

func (self Tile) PathNeighborCost(to astar.Pather) float64 {
	return 1.0
}

func (self Tile) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(Tile)
	return float64(self.Position.ManhattanDistance(toT.Position))
}

func (tile Tile) ClosestAvailableAround(to Tile) Tile {
	directions := [][2]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	}

	var queue = []Tile{to}
	var seen []Tile

	for {
		var considered []Tile
		var foundFree bool

		for _, tile := range queue {
			for _, dir := range directions {
				position := Position{X: tile.Position.X + dir[0], Y: tile.Position.Y + dir[1]}
				target := tile.env.At(position)
				if target != nil {
					considered = append(considered, *target)
					foundFree = foundFree || target.Free()
				}
			}
		}

		if foundFree {
			freeTiles := funk.Filter(considered, func(it interface{}) bool {
				return it.(Tile).Free() || it.(Tile) == to
			}).([]Tile)
			sort.Sort(byDistance{from: to, tiles: freeTiles})
			return freeTiles[0]
		}

		leftDifference, _ := funk.Difference(considered, seen)
		newQueue := leftDifference.([]interface{})
		queue = make([]Tile, 0, len(newQueue))
		for _, it := range newQueue {
			queue = append(queue, it.(Tile))
		}
	}
}
