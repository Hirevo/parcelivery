package main

import (
	"sort"

	funk "github.com/thoas/go-funk"
)

type (
	Truck struct {
		Position  Position
		Turnover  int
		MaxCharge int
		CurCharge int
		State     TruckState
		GoneUntil int
	}

	TruckState  int
)

const (
	TSWaiting TruckState = iota
	TSGone
)

func (truck *Truck) NextTurn(env *Environment) (Event, EndState) {
	undeliveredParcels := funk.Filter(env.Parcels, func(it interface{}) bool {
		return it.(*Parcel).State != PSDelivered
	}).([]*Parcel)

	sort.Sort(byWeight(undeliveredParcels))

	if len(undeliveredParcels) == 0 && truck.CurCharge == 0 {
		// all parcels have been delivered, task is gone
		return Event{}, ESCompleted
	}

	switch truck.State {
	case TSWaiting:
		if (env.TurnCount == 1) || len(undeliveredParcels) == 0 || int(undeliveredParcels[0].Weight) > (truck.MaxCharge - truck.CurCharge) {
			return Event{Kind: EKTruckGone, Truck: truck}, ESContinue
		}
		return Event{Kind: EKTruckWaiting, Truck: truck}, ESContinue
	case TSGone:
		// truck.GoneUntil--
		if truck.GoneUntil == 0 {
			// truck.State = TSWaiting
			// truck.CurCharge = 0
			return Event{Kind: EKTruckWaiting, Truck: truck}, ESContinue
		}
		return Event{Kind: EKTruckGone, Truck: truck}, ESContinue
	}

	return Event{}, ESCantMakeProgress
}
