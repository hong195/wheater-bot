package webapi

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/hong195/wheater-bot/internal/entity"
	"net/http"
	"net/url"
)

const reverseGeocodingUrl = "https://nominatim.openstreetmap.org/reverse"

type addressApiResponse struct {
	Address struct {
		County string `json:"county"`
		City   string `json:"city"`
	} `json:"address"`
}

type CityDetailsRepository struct{}

func NewCityDetailsRepository() *CityDetailsRepository {
	return &CityDetailsRepository{}
}

func (c *CityDetailsRepository) GetCityDetailsByCoordinates(ctx context.Context, lat, lon string) (*entity.CityDetails, error) {

	u, err := url.Parse(reverseGeocodingUrl)

	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("format", "json")
	q.Set("lat", lat)
	q.Set("lon", lon)
	u.RawQuery = q.Encode()

	response, err := http.Get(u.String())

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("weather API returned non-OK status: %d", response.StatusCode)
	}

	var apiResponse addressApiResponse

	if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode weather API response: %w", err)
	}

	return &entity.CityDetails{
		City:    apiResponse.Address.City,
		Country: apiResponse.Address.County,
	}, nil
}
