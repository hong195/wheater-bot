package webapi

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/hong195/wheater-bot/internal/entity"
	"io"
	"net/http"
	"net/url"
	"time"
)

type WeatherWebApi struct {
	ApiKey string
	ApiUrl string
}

func NewWeatherWebApi(apiKey, baseUrl string) *WeatherWebApi {
	return &WeatherWebApi{
		ApiKey: apiKey,
		ApiUrl: baseUrl,
	}
}

type weatherWebApiResponse struct {
	Current struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
		WindSpeed float64 `json:"wind_speed"`
		Weather   []struct {
			Main        string `json:"main"`
			Description string `json:"description"`
		} `json:"weather"`
	} `json:"current"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (w *WeatherWebApi) GetWeatherByCoordinates(ctx context.Context, lat, lon string) (*entity.Weather, error) {

	u, err := url.Parse(w.ApiUrl)

	if err != nil {
		return nil, fmt.Errorf("failed to parse weather API URL: %w", err)
	}

	q := u.Query()
	q.Set("lat", lat)
	q.Set("lon", lon)
	q.Set("appid", w.ApiKey)
	q.Set("dt", fmt.Sprintf("%d", time.Now().Unix())) // Current time for historical data
	q.Set("units", "metric")
	q.Set("lang", "ru")
	u.RawQuery = q.Encode()

	response, err := http.Get(u.String())

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("weather API returned non-OK status: %d, body: %s", response.StatusCode, string(body))
	}

	var apiResponse weatherWebApiResponse

	if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode weather API response: %w", err)
	}

	return &entity.Weather{}, nil
}
