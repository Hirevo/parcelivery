package main

import "fmt"

// DisplayStandard displays the turns of the simulation in format expected by the subject
func (env *Environment) DisplayStandard(turn int, truckEvent Event, transportEvents []Event) {
	fmt.Printf("Turn %d:\n", turn)
	for _, event := range transportEvents {
		event.Display()
	}
	truckEvent.Display()
}

// DisplayFancy displays the turns of the simulation in a fancier and more visual format
func (env *Environment) DisplayFancy(turn int, truckEvent Event, transportEvents []Event) {
	fmt.Printf("Turn %d:\n", turn)
	for curY := 0; curY < env.Size[1]; curY++ {
		for curX := 0; curX < env.Size[0]; curX++ {
			position := Position{X: curX, Y: curY}
			tile := env.At(position)
			if position == env.Truck.Position && env.Truck.State == TSWaiting {
				fmt.Print("C")
			} else {
				switch tile.Kind {
				case TKFree:
					fmt.Print(".")
				case TKParcel:
					fmt.Print("P")
				case TKTransport:
					fmt.Print("T")
				}
			}
		}
		fmt.Println()
	}
}
