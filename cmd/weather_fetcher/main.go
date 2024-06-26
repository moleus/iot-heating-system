package main

import (
	"context"
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	weatherfetcher "iot-heating-system/cmd/weather_fetcher/api"
	"iot-heating-system/pkg/common"
	"log"
	"math/rand/v2"
	"time"
)

const fetchIntervalSecs = 30

var _ = generateRandomWeather()

func init() {
	pflag.String("mqtt_broker", "tcp://mosquitto:1883", "MQTT broker to connect to")
	pflag.String("mqtt_topic", common.WeatherPredictionsTopic, "MQTT topic to write to")
	pflag.String("weather_api_url", "https://api.weather.yandex.ru/graphql/query", "Url to fetch weather from")
	pflag.String("api_key", "", "API key used for authentication")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
}

func main() {
	router := gin.Default()

	server := NewWeatherFetcher(viper.GetString("weather_api_url"), viper.GetString("api_key"))
	strictHandler := weatherfetcher.NewStrictHandler(server, nil)

	weatherfetcher.RegisterHandlers(router, strictHandler)

	go backgroundWeatherFetcher(server)

	// Start serving traffic
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// runs infinite loop and generates random weather data every 5 seconds and pushes them to mqtt
func backgroundWeatherFetcher(server *WeatherFetcher) {
	mqttClient := mqtt.NewClient(mqtt.NewClientOptions().AddBroker(viper.GetString("mqtt_broker")))
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to connect to mqtt broker: %v", token.Error())
	}

	for {
		time.Sleep(fetchIntervalSecs * time.Second)
		//43.139873051972195, 131.92746165215814 - Vladivostok
		params := weatherfetcher.GetWeatherParams{
			Lat:  43.139873051972195,
			Lon:  131.92746165215814,
			Days: 1,
		}
		weather, err := server.requestWeather(context.Background(), params)
		if err != nil {
			log.Fatalf("Failed to request weather: %v", err)
		}

		hours := weather.Data.WeatherByPoint.Forecast.Days[0].Hours

		log.Printf("Fetched weather: %v", hours)
		marshaledWeather, err := json.Marshal(hours)
		if err != nil {
			log.Fatalf("Failed to marshal weather: %v", err)
		}
		if token := mqttClient.Publish(viper.GetString("mqtt_topic"), 0, false, marshaledWeather); token.Wait() && token.Error() != nil {
			log.Fatalf("Failed to publish message: %v", token.Error())
		}
	}
}

func generateRandomWeather() []weatherfetcher.HourForecast {
	var forecasts []weatherfetcher.HourForecast

	// Get the current time and truncate to the start of the next day
	now := time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour)

	for i := 0; i < 24; i++ {
		temp := rand.Float32() * 14 // Random temperature between 0 and 14
		forecastTime := now.Add(time.Duration(i) * time.Hour)
		forecasts = append(forecasts, weatherfetcher.HourForecast{
			Temperature: temp,
			Time:        forecastTime,
		})
	}

	return forecasts
}
