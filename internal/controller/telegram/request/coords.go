package request

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	invalidCoordsErr = errors.New("incorrect coords format")
	separatorErr     = errors.New("coords should be separated by comma")
	parseErr         = errors.New("error parsing coordinates")
)

type Coords struct {
	Lat float64
	Lon float64
}

func ParseCoordinates(text string) (*Coords, error) {
	var coordText string

	if parts := strings.SplitN(text, " ", 2); len(parts) == 2 {
		coordText = parts[1]
	} else {
		coordText = text
	}

	coords := strings.Split(coordText, ",")
	if len(coords) != 2 {
		return nil, fmt.Errorf("%v", separatorErr)
	}

	lat, err1 := strconv.ParseFloat(strings.TrimSpace(coords[0]), 64)
	lon, err2 := strconv.ParseFloat(strings.TrimSpace(coords[1]), 64)

	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("%v", parseErr)
	}

	return &Coords{Lat: lat, Lon: lon}, nil
}
