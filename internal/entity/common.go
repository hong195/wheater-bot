package entity

import "errors"

var NotFoundError = errors.New("not found")

type Coords struct {
	Latitude  float64
	Longitude float64
}
