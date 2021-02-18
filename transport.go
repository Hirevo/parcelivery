package main

import (
	"sort"

	astar "github.com/beefsack/go-astar"
	funk "github.com/thoas/go-funk"
)

// Transport represents the state of a transport in the simulation
type Transport struct {
	Name     string
	Position Position
	Carrying *Parcel
}

// DeliverParcel attempts to deliver a parcel to an adjacent truck and returns whether it was successful or not
func (transport *Transport) DeliverParcel(env *Environment) (Event, bool) {
	if env.Truck.State == TSWaiting && transport.Position.IsAdjacent(env.Truck.Position) && (env.Truck.MaxCharge-env.Truck.CurCharge) >= int(transport.Carrying.Weight) {
		return Event{Kind: EKTransportLeave, Transport: transport, Parcel: transport.Carrying, Truck: env.Truck}, true
	}

	return Event{}, false
}

// PickParcel attempts to pickup an adjacent parcel and returns whether it succeeded or not
func (transport *Transport) PickParcel(env *Environment) (Event, bool) {
	x, y := transport.Position.X, transport.Position.Y

	coords := []Position{
		{X: x - 1, Y: y},
		{X: x + 1, Y: y},
		{X: x, Y: y - 1},
		{X: x, Y: y + 1},
	}

	for _, it := range coords {
		if tile := env.At(it); tile != nil {
			if tile.Kind == TKParcel {
				return Event{Kind: EKTransportTake, Transport: transport, Parcel: tile.Parcel}, true
			}
		}
	}

	return Event{}, false
}

// GoTowardsTruck makes one move closer to the truck or waits if it could not make progress reaching it
func (transport *Transport) GoTowardsTruck(env *Environment) (Event, EndState) {
	if int(transport.Carrying.Weight) > (env.Truck.MaxCharge-env.Truck.CurCharge) {
		env.NeedsToGo = true
	}

	from := *env.At(transport.Position)
	to := *env.At(env.Truck.Position)
	path, _, ok := astar.Path(from, from.ClosestAvailableAround(to))
	if !ok {
		return Event{Kind: EKTransportWaiting, Transport: transport}, ESContinue
	}
	var tiles []Tile
	for _, it := range path {
		tiles = append(tiles, it.(Tile))
	}
	if len(tiles) < 2 {
		return Event{Kind: EKTransportWaiting, Transport: transport}, ESContinue
	}

	moveTo := tiles[len(tiles)-2].Position
	return Event{Kind: EKTransportGo, Transport: transport, To: moveTo}, ESContinue
}

// GoTowardsClosestParcel makes one move closer to the the nearest parcel or waits if it could not make progress reaching it
func (transport *Transport) GoTowardsClosestParcel(env *Environment) (Event, EndState) {
	from := *env.At(transport.Position)

	freeParcels := funk.Filter(env.Parcels, func(it interface{}) bool {
		parcel := it.(*Parcel)
		return parcel.State == PSFree
	}).([]*Parcel)

	if len(freeParcels) == 0 {
		// got nothing to do
		return Event{Kind: EKTransportWaiting, Transport: transport}, ESContinue
	}

	deliverableParcels := funk.Filter(freeParcels, func(it interface{}) bool {
		parcel := it.(*Parcel)
		return int(parcel.Weight) <= (env.Truck.MaxCharge - env.Truck.CurCharge)
	}).([]*Parcel)

	if len(deliverableParcels) == 0 {
		// nothing atteignable is deliverable, pinging the truck it may need to go
		env.NeedsToGo = true
		deliverableParcels = freeParcels
	}

	var paths [][]Tile
	for _, parcel := range deliverableParcels {
		to := *env.At(parcel.Position)
		path, _, ok := astar.Path(from, to.ClosestAvailableAround(from))
		if ok && len(path) > 1 {
			var tiles []Tile
			for _, it := range path {
				tiles = append(tiles, it.(Tile))
			}
			paths = append(paths, tiles)
		}
	}

	sort.Sort(byPathLength(paths))
	if len(paths) == 0 {
		// no paths found.
		return Event{Kind: EKTransportWaiting, Transport: transport}, ESContinue
	}

	moveTo := paths[0][len(paths[0])-2].Position
	return Event{Kind: EKTransportGo, Transport: transport, To: moveTo}, ESContinue
}

// NextTurn computes what this transport should do for the turn
func (transport *Transport) NextTurn(env *Environment) (Event, EndState) {
	if transport.Carrying != nil {
		if event, done := transport.DeliverParcel(env); done {
			return event, ESContinue
		}
		return transport.GoTowardsTruck(env)
	}

	if event, done := transport.PickParcel(env); done {
		return event, ESContinue
	}

	return transport.GoTowardsClosestParcel(env)
}
