package webapi

import (
	"context"
	"errors"
	"github.com/hong195/wheater-bot/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"strings"
	"testing"
)

const testApiUrl = "https://api.test.com"
const testApiKey = "test_api_key"
const invalidApiKey = "invalid_api_key"

func TestGetWeatherDetails(t *testing.T) {

	mockTransport, client := getMockHttpClient()

	repo := NewWeatherWebApi(client, testApiKey, testApiUrl)
	body := `{"current": {"sunrise": 1718160000,"sunset": 1718203200,"temp": 27.5,"feels_like": 29.1,"pressure": 1013,"humidity": 60,"visibility": 10000,"wind_speed": 3.2,"wind_deg": 180,"wind_gust": 5.6,"weather": [{"id": 801,"main": "Clouds","description": "few clouds"}]}}`

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}

	mockTransport.On("RoundTrip", mock.Anything).Return(response, nil)

	res, _ := repo.GetWeatherByCoordinates(context.TODO(), testLat, testLon)
	exp := &entity.Weather{
		Sunrise:     1718160000,
		Sunset:      1718203200,
		Temperature: 27.5,
		FeelsLike:   29.1,
		WindSpeed:   3.2,
		Humidity:    60,
		Description: "few clouds",
	}

	assert.NotNil(t, repo, "repo cannot be nil!")
	assert.Equal(t, res, exp, "error in fetching weather details")
}

func TestFailsIfApiUrlIsInvalid(t *testing.T) {
	_, client := getMockHttpClient()

	invalidUrl := "$234234_3%"
	repo := NewWeatherWebApi(client, testApiKey, invalidUrl)

	_, err := repo.GetWeatherByCoordinates(context.TODO(), testLat, testLon)

	assert.Error(t, err, "api url is invalid")
}

func TestFailsIfApiKeyIsInvalid(t *testing.T) {
	mockTransport, client := getMockHttpClient()

	repo := NewWeatherWebApi(client, invalidApiKey, testApiKey)

	response := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
	}

	mockTransport.On("RoundTrip", mock.Anything).Return(response, errors.New("api url is invalid"))
	_, err := repo.GetWeatherByCoordinates(context.TODO(), testLat, testLon)

	assert.Error(t, err, "api url is invalid")
}

func TestFailsIfReturnNot200Ok(t *testing.T) {
	mockTransport, client := getMockHttpClient()

	repo := NewWeatherWebApi(client, invalidApiKey, testApiKey)

	response := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(strings.NewReader("")),
		Header:     make(http.Header),
	}

	mockTransport.On("RoundTrip", mock.Anything).Return(response, nil)
	_, err := repo.GetWeatherByCoordinates(context.TODO(), testLat, testLon)

	assert.Error(t, err, "bad request")
}

func TestFailsIfJsonIsBroken(t *testing.T) {

	mockTransport, client := getMockHttpClient()

	repo := NewWeatherWebApi(client, testApiKey, testApiUrl)
	invalidBody := `{"current": {"sunrise": 1718160000,"sunset"`

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(invalidBody)),
		Header:     make(http.Header),
	}

	mockTransport.On("RoundTrip", mock.Anything).Return(response, nil)

	_, err := repo.GetWeatherByCoordinates(context.TODO(), testLat, testLon)

	assert.Error(t, err, "invalid json parsing")
}

func TestDescriptionMightBeEmpty(t *testing.T) {
	mockTransport, client := getMockHttpClient()

	repo := NewWeatherWebApi(client, testApiKey, testApiUrl)
	body := `{"current": {"sunrise": 1718160000,"sunset": 1718203200,"temp": 27.5,"feels_like": 29.1,"pressure": 1013,"humidity": 60,"visibility": 10000,"wind_speed": 3.2,"wind_deg": 180,"wind_gust": 5.6,"weather": [{"id": 801,"main": "Clouds","description": ""}]}}`

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}

	mockTransport.On("RoundTrip", mock.Anything).Return(response, nil)

	res, err := repo.GetWeatherByCoordinates(context.TODO(), testLat, testLon)

	assert.Nil(t, err, "error should be nil")

	assert.Empty(t, res.Description, "description should be empty")
}
