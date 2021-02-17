package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	funk "github.com/thoas/go-funk"
)

func parseFromFile(path string) (*Environment, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("input: could not open file: %w", err)
	}

    env := &Environment{}

    contents := string(bytes)
    // trimming allows to remove a potential trailing newline
    contents = strings.TrimSpace(contents)

	lines := strings.Split(contents, "\n")
	lineCount := len(lines)

	if len(lines) < 2 {
		return nil, fmt.Errorf("input: file comports less than the minimum of 2 lines")
	}

	initial := lines[0]
	initialParts := strings.Split(initial, " ")
	if len(initialParts) != 3 {
        return nil, fmt.Errorf("input(line 1): unexpected count of components for environment")
	}
	if env.Size[0], err = strconv.Atoi(initialParts[0]); err != nil {
		return nil, fmt.Errorf("input(line 1): could not parse the map width: %w", err)
	}
	if env.Size[1], err = strconv.Atoi(initialParts[1]); err != nil {
		return nil, fmt.Errorf("input(line 1): could not parse the map height: %w", err)
	}
	if env.TurnCount, err = strconv.Atoi(initialParts[2]); err != nil {
		return nil, fmt.Errorf("input(line 1): could not parse the turn count: %w", err)
	}

	currentLine := 1
	for currentLine < (lineCount - 1) {
		parcelParts := strings.Split(lines[currentLine], " ")
		if len(parcelParts) == 3 {
			// this line is potentially describing a transport, moving onto the next stage...
			break
		}
		parcel := Parcel{}
		if len(parcelParts) != 4 {
			return nil, fmt.Errorf("input(line %d): unexpected count of components for parcel", currentLine+1)
		}
		parcel.Name = parcelParts[0]
		if len(parcel.Name) == 0 {
			return nil, fmt.Errorf("input(line %d): empty name for parcel", currentLine+1)
		}
		if parcel.Position.X, err = strconv.Atoi(parcelParts[1]); err != nil {
			return nil, fmt.Errorf("input(line %d): could not parse X position for parcel: %w", currentLine+1, err)
		}
		if parcel.Position.Y, err = strconv.Atoi(parcelParts[2]); err != nil {
			return nil, fmt.Errorf("input(line %d): could not parse Y position for parcel: %w", currentLine+1, err)
		}
		switch strings.ToLower(parcelParts[3]) {
		case "green":
			parcel.Weight = WGHTGreen
		case "blue":
			parcel.Weight = WGHTBlue
		case "yellow":
			parcel.Weight = WGHTYellow
		default:
			return nil, fmt.Errorf("input(line %d): could not parse weight `%s` for parcel", currentLine+1, parcelParts[3])
		}
		env.Parcels = append(env.Parcels, &parcel)
		currentLine++
	}

	for currentLine < (lineCount - 1) {
		transportParts := strings.Split(lines[currentLine], " ")
		if len(transportParts) != 3 {
			return nil, fmt.Errorf("input(line %d): unexpected count of components for transport", currentLine+1)
		}
		transport := Transport{}
        transport.Name = transportParts[0]
		if len(transport.Name) == 0 {
			return nil, fmt.Errorf("input(line %d): empty name for transport", currentLine+1)
		}
		if transport.Position.X, err = strconv.Atoi(transportParts[1]); err != nil {
			return nil, fmt.Errorf("input(line %d): could not parse X position for transport: %w", currentLine+1, err)
		}
		if transport.Position.Y, err = strconv.Atoi(transportParts[2]); err != nil {
			return nil, fmt.Errorf("input(line %d): could not parse Y position for transport: %w", currentLine+1, err)
		}
		env.Transports = append(env.Transports, &transport)
		currentLine++
    }

    truck := Truck{}
    truckLine := lines[currentLine]
	truckParts := strings.Split(truckLine, " ")
	if len(truckParts) != 4 {
        return nil, fmt.Errorf("input(line %d): unexpected count of components for truck", currentLine+1)
	}
	if truck.Position.X, err = strconv.Atoi(truckParts[0]); err != nil {
		return nil, fmt.Errorf("input(line %d): could not parse X position for truck: %w", currentLine+1, err)
	}
	if truck.Position.Y, err = strconv.Atoi(truckParts[1]); err != nil {
		return nil, fmt.Errorf("input(line %d): could not parse Y position for truck: %w", currentLine+1, err)
	}
	if truck.MaxCharge, err = strconv.Atoi(truckParts[2]); err != nil {
		return nil, fmt.Errorf("input(line %d): could not parse maximum charge for truck: %w", currentLine+1, err)
	}
	if truck.Turnover, err = strconv.Atoi(truckParts[3]); err != nil {
		return nil, fmt.Errorf("input(line %d): could not parse turnover delay for truck: %w", currentLine+1, err)
    }
    env.Truck = &truck

	transportNames := funk.Map(env.Transports, func(it interface{}) interface{} {
		return it.(*Transport).Name
	}).([]interface{})
	parcelNames := funk.Map(env.Parcels, func(it interface{}) interface{} {
		return it.(*Parcel).Name
	}).([]interface{})

	uniqueNames := funk.Uniq(funk.FlattenDeep([][]interface{}{transportNames, parcelNames})).([]interface{})

	if len(uniqueNames) != len(transportNames) + len(parcelNames) {
		return nil, fmt.Errorf("input: some parcel(s) and/or transport(s) share common names")
	}

	return env, nil
}
