package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	weatherfetcher "iot-heating-system/cmd/weather_fetcher/api"
	"log"
	"net/http"
	"strings"
)

var _ weatherfetcher.StrictServerInterface = (*WeatherFetcher)(nil)

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
				weatherfetcher.DaysForecast
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

func (w *WeatherFetcher) GetWeather(c context.Context, req weatherfetcher.GetWeatherRequestObject) (weatherfetcher.GetWeatherResponseObject, error) {
	yandexResponse, err := w.requestWeather(c, req.Params)
	if err != nil {
		return weatherfetcher.GetWeather500Response{}, errors.Wrap(err, "failed to request weather from Yandex API")
	}

	return weatherfetcher.GetWeather200JSONResponse{Days: yandexResponse.Data.WeatherByPoint.Forecast.Days}, nil
}

func (w *WeatherFetcher) requestWeather(c context.Context, req weatherfetcher.GetWeatherParams) (YandexResponse, error) {
	graphQLBody := fmt.Sprintf(yandexForecastGraphqlTemplate, req.Lat, req.Lon, req.Days)
	graphQLBody = strings.ReplaceAll(graphQLBody, "\t", "")

	r, err := http.NewRequestWithContext(c, http.MethodPost, w.WeatherApiUrl, strings.NewReader(graphQLBody))
	if err != nil {
		return YandexResponse{}, errors.Wrap(err, "create request")
	}
	r.Header.Add("X-Yandex-Weather-Key", w.ApiKey)
	r.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return YandexResponse{}, errors.Wrap(err, "request Yandex API")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return YandexResponse{}, errors.Wrap(err, "status code not 200 from Yandex API")
	}

	// Assuming you want to return the response from the weather API directly
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return YandexResponse{}, errors.Wrap(err, "read response body from Yandex API")
	}

	log.Printf("Got response from Weather API: %s", body)

	yandexResponse := YandexResponse{}

	err = json.Unmarshal(body, &yandexResponse)
	if err != nil {
		return YandexResponse{}, errors.Wrap(err, "unmarshal body from Yandex API")
	}
	return yandexResponse, nil
}
