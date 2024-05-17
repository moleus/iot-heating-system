package main

import (
	"github.com/gin-gonic/gin"
	"iot-heating-system/api"
	"time"
)

var _ api.ServerInterface = (*WeatherFetcher)(nil)

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

func (w *WeatherFetcher) GetWeather(c *gin.Context, params api.GetWeatherParams) {
	mockWeather := w.getMockWeather()
	c.JSON(200, mockWeather)
}

func (w *WeatherFetcher) getMockWeather() api.DaysForecast {
	hourForecasts := []api.HourForecast{
		{
			Temperature: new(float32),
			Time:        new(time.Time),
		},
		{
			Temperature: new(float32),
			Time:        new(time.Time),
		},
		{
			Temperature: new(float32),
			Time:        new(time.Time),
		},
	}
	dayForecasts := []api.DayForecast{
		{
			Forecasts: &hourForecasts,
		},
		{
			Forecasts: &hourForecasts,
		},
	}
	return dayForecasts
}
