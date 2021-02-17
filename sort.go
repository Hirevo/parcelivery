package main

type byWeight []*Parcel

func (a byWeight) Len() int           { return len(a) }
func (a byWeight) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byWeight) Less(i, j int) bool { return a[i].Weight < a[j].Weight }

type byPathLength [][]Tile

func (a byPathLength) Len() int           { return len(a) }
func (a byPathLength) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPathLength) Less(i, j int) bool { return len(a[i]) < len(a[j]) }

type byDistance struct {
	tiles []Tile
	from  Tile
}

func (a byDistance) Len() int      { return len(a.tiles) }
func (a byDistance) Swap(i, j int) { a.tiles[i], a.tiles[j] = a.tiles[j], a.tiles[i] }
func (a byDistance) Less(i, j int) bool {
	return a.tiles[i].PathEstimatedCost(a.from) < a.tiles[j].PathEstimatedCost(a.from)
}
