package main

import "fmt"

type (
	// Event represents an event in the simulation, like a truck move or a parcel deposit
	Event struct {
		Kind      EventKind
		Transport *Transport
		Truck     *Truck
		To        Position
		Parcel    *Parcel
	}

	// EventKind is the discriminant for `Event`
	EventKind int
)

const (
	EKTruckWaiting EventKind = iota
	EKTruckGone
	EKTransportWaiting
	EKTransportGo
	EKTransportTake
	EKTransportLeave
)

func (event *Event) Display() {
	switch event.Kind {
	case EKTruckWaiting:
		fmt.Printf("camion WAITING %d %d\n", event.Truck.MaxCharge, event.Truck.CurCharge)
	case EKTruckGone:
		fmt.Printf("camion GONE %d %d\n", event.Truck.MaxCharge, event.Truck.CurCharge)
	case EKTransportWaiting:
		fmt.Printf("%s WAIT\n", event.Transport.Name)
	case EKTransportGo:
		fmt.Printf("%s GO [%d,%d]\n", event.Transport.Name, event.Transport.Position.X, event.Transport.Position.Y)
	case EKTransportTake:
		fmt.Printf("%s TAKE %s %s\n", event.Transport.Name, event.Parcel.Name, WeightToColor(event.Parcel.Weight))
	case EKTransportLeave:
		fmt.Printf("%s LEAVE %s %s\n", event.Transport.Name, event.Parcel.Name, WeightToColor(event.Parcel.Weight))
	}
}
