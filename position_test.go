package main

import (
	"fmt"
	"testing"
)

func TestManhattanDistance(t *testing.T) {
	table := []struct {
		input    [2]Position
		expected int
	}{
		{input: [2]Position{{X: 0, Y: 0}, {X: 0, Y: 0}}, expected: 0},
		{input: [2]Position{{X: 0, Y: 0}, {X: 1, Y: 0}}, expected: 1},
		{input: [2]Position{{X: 0, Y: 0}, {X: 1, Y: 1}}, expected: 2},
		{input: [2]Position{{X: 1, Y: 0}, {X: 0, Y: 0}}, expected: 1},
		{input: [2]Position{{X: 1, Y: 1}, {X: 0, Y: 0}}, expected: 2},
		{input: [2]Position{{X: 2, Y: 5}, {X: 3, Y: 2}}, expected: 4},
		{input: [2]Position{{X: 6, Y: 2}, {X: 2, Y: 6}}, expected: 8},
	}
	for idx, tt := range table {
		t.Run(fmt.Sprintf("TestManhattanDistance.%d", idx), func(t *testing.T) {
			s := tt.input[0].ManhattanDistance(tt.input[1])
			if s != tt.expected {
				t.Errorf("for input \"%+v\", got %v, expected %v", tt.input, s, tt.expected)
			}
		})
	}
}
