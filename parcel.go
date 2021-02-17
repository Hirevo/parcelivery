package main

type (
	Parcel struct {
		Name     string
		Position Position
		Weight   Weight
		State    ParcelState
	}

	Weight      int
	ParcelState int
)

const (
	WGHTYellow Weight = 100
	WGHTGreen         = 200
	WGHTBlue          = 500
)

const (
	PSFree ParcelState = iota
	PSCarried
	PSDelivered
)

func WeightToColor(weight Weight) string {
	switch weight {
	case WGHTYellow:
		return "YELLOW"
	case WGHTGreen:
		return "GREEN"
	case WGHTBlue:
		return "BLUE"
	}
	return "??"
}
