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

const testLat = 55.751244
const testLon = 37.618423
const testUrl = "http://localhost:8080"

type MockRoundTripper struct {
	mock.Mock
}

func (m *MockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	args := m.Called(req)

	resp := args.Get(0).(*http.Response)
	return resp, args.Error(1)
}

func getMockHttpClient() (*MockRoundTripper, *http.Client) {
	mockTransport := new(MockRoundTripper)
	client := &http.Client{
		Transport: mockTransport,
	}

	return mockTransport, client
}

func TestCityRepoCanFetchDetailsByLatLon(t *testing.T) {
	mockTransport, client := getMockHttpClient()

	repo := NewCityDetailsRepository(client, testUrl)
	body := `{"address": {"county": "Uzbekistan","city": "Tashkent"}}`

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}

	mockTransport.On("RoundTrip", mock.Anything).Return(response, nil)

	res, _ := repo.GetCityDetailsByCoordinates(context.TODO(), testLat, testLon)
	exp := &entity.CityDetails{
		City:    "Tashkent",
		Country: "Uzbekistan",
	}

	assert.NotNil(t, repo)
	assert.Equal(t, res, exp)
}

func TestRepoReturnsErrorIfApiUrlIsInvalid(t *testing.T) {
	_, client := getMockHttpClient()

	invalidUrl := "$234234_3%"
	repo := NewCityDetailsRepository(client, invalidUrl)

	_, err := repo.GetCityDetailsByCoordinates(context.TODO(), testLat, testLon)

	assert.Error(t, err)
}

func TestCantHandleFailedApiInnerRequest(t *testing.T) {
	mockTransport, client := getMockHttpClient()

	repo := NewCityDetailsRepository(client, testUrl)

	response := &http.Response{}

	mockTransport.On("RoundTrip", mock.Anything).Return(response, errors.New("some error"))

	_, err := repo.GetCityDetailsByCoordinates(context.TODO(), testLat, testLon)

	assert.Error(t, err)
}

func TestRepoInnerApiCallAccepts200Status(t *testing.T) {
	mockTransport, client := getMockHttpClient()

	repo := NewCityDetailsRepository(client, testUrl)

	response := &http.Response{
		StatusCode: http.StatusInternalServerError,
	}

	mockTransport.On("RoundTrip", mock.Anything).Return(response, nil)

	_, err := repo.GetCityDetailsByCoordinates(context.TODO(), testLat, testLon)

	assert.Error(t, err, http.StatusInternalServerError)
}

func TestGetCityDetailsByCoordinatesFailsIfCantParseJson(t *testing.T) {
	mockTransport, client := getMockHttpClient()
	repo := NewCityDetailsRepository(client, testUrl)

	invalidInnerApiRes := `{"address": {"county": "Uzbekistan","city": }`

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(invalidInnerApiRes)),
		Header:     make(http.Header),
	}

	mockTransport.On("RoundTrip", mock.Anything).Return(response, nil)

	_, err := repo.GetCityDetailsByCoordinates(context.TODO(), testLat, testLon)

	assert.Error(t, err, "failed to decode weather API response")
}
