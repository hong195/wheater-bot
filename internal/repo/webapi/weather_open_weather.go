package webapi

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/hong195/wheater-bot/internal/entity"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type WeatherWebApi struct {
	httpClient *http.Client
	ApiKey     string
	ApiUrl     string
}

func NewWeatherWebApi(httpClient *http.Client, apiKey, baseUrl string) *WeatherWebApi {
	return &WeatherWebApi{
		httpClient: httpClient,
		ApiKey:     apiKey,
		ApiUrl:     baseUrl,
	}
}

type WeatherWebApiResponse struct {
	Current struct {
		Sunrise    int     `json:"sunrise"`
		Sunset     int     `json:"sunset"`
		Temp       float64 `json:"temp"`
		FeelsLike  float64 `json:"feels_like"`
		Pressure   int     `json:"pressure"`
		Humidity   int     `json:"humidity"`
		Visibility int     `json:"visibility"`
		WindSpeed  float64 `json:"wind_speed"`
		WindDeg    int     `json:"wind_deg"`
		WindGust   float64 `json:"wind_gust"`
		Weather    []struct {
			Description string `json:"description"`
		}
	} `json:"current"`
}

func (w *WeatherWebApi) GetWeatherByCoordinates(ctx context.Context, lat, lon float64) (*entity.Weather, error) {

	u, err := url.Parse(w.ApiUrl)
	currentTime := time.Now().Unix()

	if err != nil {
		return nil, fmt.Errorf("failed to parse weather API URL: %w", err)
	}

	q := u.Query()
	q.Set("format", "json")
	q.Set("lat", strconv.FormatFloat(lat, 'f', -1, 64))
	q.Set("lon", strconv.FormatFloat(lon, 'f', -1, 64))
	q.Set("appid", w.ApiKey)
	q.Set("dt", fmt.Sprintf("%d", currentTime))
	q.Set("units", "metric")
	q.Set("lang", "ru")
	u.RawQuery = q.Encode()

	response, err := w.httpClient.Get(u.String())

	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("weather API returned non-OK status: %d, body: %s", response.StatusCode, string(body))
	}

	var apiResponse WeatherWebApiResponse

	if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode weather API response: %w", err)
	}

	desc := ""

	if apiResponse.Current.Weather[0].Description != "" {
		desc = apiResponse.Current.Weather[0].Description
	}

	return &entity.Weather{
		Sunrise:     apiResponse.Current.Sunrise,
		Sunset:      apiResponse.Current.Sunset,
		Temperature: apiResponse.Current.Temp,
		Humidity:    apiResponse.Current.Humidity,
		FeelsLike:   apiResponse.Current.FeelsLike,
		WindSpeed:   apiResponse.Current.WindSpeed,
		Description: desc,
	}, nil
}
