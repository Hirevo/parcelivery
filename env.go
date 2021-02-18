package main

type (
	Environment struct {
		Size       [2]int
		TurnCount  int
		Parcels    []*Parcel
		Transports []*Transport
		Truck      *Truck

		// used to ping the truck that it needs to go to empty its charge
		NeedsToGo  bool
	}

	EndState int
)

const (
	ESContinue EndState = iota
	ESCompleted
	ESTurnCountExpired
	ESCantMakeProgress
)

// At gets the tile data at a given position
func (env *Environment) At(pos Position) *Tile {
	if pos.X < 0 || pos.X >= env.Size[0] || pos.Y < 0 || pos.Y >= env.Size[1] {
		return nil
	}

	for _, it := range env.Parcels {
		if pos == it.Position && it.State == PSFree {
			return &Tile{env: env, Kind: TKParcel, Position: pos, Parcel: it}
		}
	}
	for _, it := range env.Transports {
		if pos == it.Position {
			return &Tile{env: env, Kind: TKTransport, Position: pos, Transport: it}
		}
	}

	return &Tile{env: env, Kind: TKFree, Position: pos}
}

// NextTurn computes all the moves within an environment for a turn
func (env *Environment) NextTurn() (Event, []Event, EndState) {
	var truckEvent Event
	var transportEvents []Event

	if env.TurnCount == 0 {
		return truckEvent, transportEvents, ESTurnCountExpired
	}
	env.TurnCount--

	if event, state := env.Truck.NextTurn(env); state != ESContinue {
		return truckEvent, transportEvents, state
	} else {
		env.ApplyEvent(event)
		truckEvent = event
	}

	for _, it := range env.Transports {
		if event, state := it.NextTurn(env); state != ESContinue {
			return truckEvent, transportEvents, state
		} else {
			env.ApplyEvent(event)
			transportEvents = append(transportEvents, event)
		}
	}

	return truckEvent, transportEvents, ESContinue
}

// ApplyEvent applies the consequences of an event onto the environment
func (env *Environment) ApplyEvent(event Event) {
	switch event.Kind {
	case EKTruckWaiting:
		if event.Truck.State == TSGone {
			event.Truck.GoneUntil = 0
		}
		event.Truck.State = TSWaiting
	case EKTruckGone:
		if event.Truck.State == TSWaiting {
			event.Truck.CurCharge = 0
			event.Truck.GoneUntil = event.Truck.Turnover
		}
		event.Truck.State = TSGone
		event.Truck.GoneUntil--
	case EKTransportWaiting:
	case EKTransportGo:
		event.Transport.Position = event.To
	case EKTransportTake:
		event.Parcel.State = PSCarried
		event.Transport.Carrying = event.Parcel
	case EKTransportLeave:
		event.Parcel.State = PSDelivered
		event.Transport.Carrying = nil
		event.Truck.CurCharge += int(event.Parcel.Weight)
	}
}
