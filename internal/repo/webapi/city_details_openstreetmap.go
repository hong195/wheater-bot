package webapi

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/hong195/wheater-bot/internal/entity"
	"net/http"
	"net/url"
	"strconv"
)

type addressApiResponse struct {
	Address struct {
		County       string `json:"county"`
		Country      string `json:"country"`
		City         string `json:"city"`
		PostCode     string `json:"postcode"`
		Neighborhood string `json:"neighbourhood"`
	} `json:"address"`
}

type CityDetailsRepository struct {
	httpClient *http.Client
	apiUrl     string
}

func NewCityDetailsRepository(httpClient *http.Client, apiUrl string) *CityDetailsRepository {
	return &CityDetailsRepository{httpClient, apiUrl}
}

func (c *CityDetailsRepository) GetCityDetailsByCoordinates(ctx context.Context, lat, lon float64) (*entity.CityDetails, error) {

	u, err := url.Parse(c.apiUrl)

	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("format", "jsonv2")
	q.Set("lat", strconv.FormatFloat(lat, 'f', -1, 64))
	q.Set("lon", strconv.FormatFloat(lon, 'f', -1, 64))
	q.Set("accept-language", "ru")
	u.RawQuery = q.Encode()

	response, err := c.httpClient.Get(u.String())

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
		City:          apiResponse.Address.City,
		Country:       apiResponse.Address.Country,
		County:        apiResponse.Address.County,
		PostCode:      apiResponse.Address.PostCode,
		Neighbourhood: apiResponse.Address.Neighborhood,
	}, nil
}
