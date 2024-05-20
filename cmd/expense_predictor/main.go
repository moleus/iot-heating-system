package main

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"io"
	fuelanalyzer "iot-heating-system/cmd/fuel_analyzer/api"
	"iot-heating-system/pkg/common"
	"log"
	"net/http"
	"net/url"
)

const (
	RequiredTemp                 = 20
	Efficiency                   = 0.92
	SpecificHeatOfCombustionFuel = 33500
	DesignOutsideTemp            = -24
	DesignCapacity               = 21.5
)

func init() {
	pflag.String("mqtt_broker", "tcp://mosquitto:1883", "MQTT broker to connect to")
	pflag.String("mqtt_subscribe_prediction_temperature_topic", "predictions/weather", "Prediction temperature reading topic")
	pflag.String("mqtt_publish_fuel_expenses_prediction_topic", "predictions/fuel_expenses", "Fuel expenses prediction writing topic")
	pflag.String("fuel_analyzer_url", "http://fuel-analyzer:8080/analyze", "Fuel analyzer URL to send GET request to")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
}

func main() {
	opts := mqtt.NewClientOptions().AddBroker(viper.GetString("mqtt_broker"))
	mqttClient := mqtt.NewClient(opts)

	token := mqttClient.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("Failed to connect to mqtt broker: %v", token.Error())
	}

	token = mqttClient.Subscribe(viper.GetString("mqtt_subscribe_temperature_topic"), 0, onTemperatureChange)
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %v", token.Error())
	}

	select {}
}

func onTemperatureChange(client mqtt.Client, msg mqtt.Message) {
	var hourForecast []common.HourForecast

	err := json.Unmarshal(msg.Payload(), &hourForecast)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	u, err := url.Parse(viper.GetString("fuel_analyzer_url"))
	if err != nil {
		log.Printf("Failed to parse URL: %v", err)
		return
	}

	var fuelExpensesPrediction common.MqttTargetFuelExpensesPredictions

	for _, forecast := range hourForecast {
		q := u.Query()
		q.Set("required_temp", fmt.Sprintf("%d", RequiredTemp))
		q.Set("outside_temp", fmt.Sprintf("%f", forecast.Temperature))
		q.Set("efficiency", fmt.Sprintf("%f", Efficiency))
		q.Set("specific_heat_of_combustion_fuel", fmt.Sprintf("%d", SpecificHeatOfCombustionFuel))
		q.Set("design_outside_temp", fmt.Sprintf("%d", DesignOutsideTemp))
		q.Set("design_capacity", fmt.Sprintf("%f", DesignCapacity))
		u.RawQuery = q.Encode()
		resp, err := http.Get(u.String())
		if err != nil {
			log.Printf("Failed to make GET request: %v", err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Failed to read response body: %v", err)
			return
		}

		var fuelConsumptionResponse fuelanalyzer.GetAnalyze200JSONResponse
		err = json.Unmarshal(body, &fuelConsumptionResponse)
		if err != nil {
			log.Printf("Failed to unmarshal response body: %v", err)
			return
		}
		fuelExpensesPrediction.Values = append(fuelExpensesPrediction.Values, *fuelConsumptionResponse.FuelConsumption)
		fuelExpensesPrediction.Time = append(fuelExpensesPrediction.Time, forecast.Time)
	}

	marshaledFuelExpensesPrediction, err := json.Marshal(fuelExpensesPrediction)
	if err != nil {
		log.Printf("Failed to marshal fuel expenses: %v", err)
		return
	}

	token := client.Publish(viper.GetString("mqtt_publish_fuel_expenses_prediction_topic"), 0, false, marshaledFuelExpensesPrediction)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Failed to publish message: %v", token.Error())
	}
}
