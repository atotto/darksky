package darksky

import (
	"fmt"
	"github.com/google/go-querystring/query"
)

const (
	baseUrl = "https://api.darksky.net/forecast"
)

// DarkSky Api client
type DarkSky interface {
	Forecast(request ForecastRequest) (ForecastResponse, error)
	SetBaseUrl(baseUrl string)
}

type darkSky struct {
	ApiKey  string
	BaseUrl string
}

// New creates a new DarkSky client
func New(apiKey string) DarkSky {
	return &darkSky{apiKey, baseUrl}
}

// Forecast request a forecast
func (d *darkSky) Forecast(request ForecastRequest) (ForecastResponse, error) {
	response := ForecastResponse{}

	url := d.buildRequestUrl(request)

	err := get(url, nil, &response)

	return response, err
}

func (d *darkSky) SetBaseUrl(baseUrl string) {
	d.BaseUrl = baseUrl
}

func (d *darkSky) buildRequestUrl(request ForecastRequest) string {
	url := fmt.Sprintf("%s/%s/%f,%f", d.BaseUrl, d.ApiKey, request.Latitude, request.Longitude)

	if request.Time != nil {
		url = url + fmt.Sprintf(",%d", request.Time)
	}

	values, _ := query.Values(request.Options)
	queryString := values.Encode()

	if len(queryString) > 0 {
		url = url + "?" + queryString
	}

	return url
}
