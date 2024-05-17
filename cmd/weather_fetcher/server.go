package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"log"
	"net/http"
	"strings"
)

var _ StrictServerInterface = (*WeatherFetcher)(nil)

const yandexForecastGraphqlTemplate = `{
	"query": "{
	  weatherByPoint(request: { lat: %f, lon: %f }) {
		forecast {
		  days(limit: %d) {
			hours {
			  time
			  temperature
			}
		  }
		}
	  }
	}"
}`

type YandexResponse struct {
	Data struct {
		WeatherByPoint struct {
			Forecast struct {
				DaysForecast
			} `json:"forecast"`
		} `json:"weatherByPoint"`
	} `json:"data"`
}

type WeatherFetcher struct {
	WeatherApiUrl string
	ApiKey        string
}

func NewWeatherFetcher(apiUrl, apiKey string) *WeatherFetcher {
	return &WeatherFetcher{
		WeatherApiUrl: apiUrl,
		ApiKey:        apiKey,
	}
}

func (w *WeatherFetcher) GetWeather(c context.Context, req GetWeatherRequestObject) (GetWeatherResponseObject, error) {
	graphQLBody := fmt.Sprintf(yandexForecastGraphqlTemplate, req.Params.Lat, req.Params.Lon, req.Params.Days)
	graphQLBody = strings.ReplaceAll(graphQLBody, "\t", "")

	r, err := http.NewRequestWithContext(c, http.MethodPost, w.WeatherApiUrl, strings.NewReader(graphQLBody))
	if err != nil {
		return GetWeather500Response{}, errors.Wrap(err, "create request")
	}
	r.Header.Add("X-Yandex-Weather-Key", w.ApiKey)
	r.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return GetWeather500Response{}, errors.Wrap(err, "request Yandex API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return GetWeather500Response{}, errors.Wrap(err, "status code not 200 from Yandex API")
	}

	// Assuming you want to return the response from the weather API directly
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return GetWeather500Response{}, errors.Wrap(err, "read response body from Yandex API")
	}

	log.Printf("Got response from Weather API: %s", body)

	yandexResponse := YandexResponse{}

	err = json.Unmarshal(body, &yandexResponse)
	if err != nil {
		return GetWeather500Response{}, errors.Wrap(err, "unmarshal body from Yandex API")
	}

	return GetWeather200JSONResponse{yandexResponse.Data.WeatherByPoint.Forecast.Days}, nil
}
