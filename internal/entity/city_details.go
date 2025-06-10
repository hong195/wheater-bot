package entity

import "context"

type CityDetailRepository interface {
	GetCityDetailsByCoordinates(ctx context.Context, lot, lat string) (*CityDetails, error)
}

type CityDetails struct {
	City    string
	Country string
}
