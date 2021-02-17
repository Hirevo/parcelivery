package main

type Position struct {
	X, Y int
}

// ManhattanDistance computes the manhattan distance between two positions
func (pos Position) ManhattanDistance(oth Position) int {
	absX := oth.X - pos.X
	if absX < 0 {
		absX = -absX
	}
	absY := oth.Y - pos.Y
	if absY < 0 {
		absY = -absY
	}
	return absX + absY
}

// IsAdjacent checks if the position is adjacent to another
func (pos Position) IsAdjacent(oth Position) bool {
	return pos.ManhattanDistance(oth) == 1
}
