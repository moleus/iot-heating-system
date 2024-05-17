package main

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

func init() {
	pflag.String("weather_api_url", "https://api.weather.yandex.ru/graphql/query", "Url to fetch weather from")
	pflag.String("api_key", "", "API key used for authentication")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
}

func main() {
	router := gin.Default()

	server := NewWeatherFetcher(viper.GetString("weather_api_url"), viper.GetString("api_key"))
	strictHandler := NewStrictHandler(server, nil)

	RegisterHandlers(router, strictHandler)

	// Start serving traffic
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
