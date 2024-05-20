package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"iot-heating-system/pkg/common"
	"log"
	"math/rand/v2"
	"time"
)

func init() {
	pflag.String("mqtt_broker", "tcp://mosquitto:1883", "MQTT broker to connect to")
	pflag.String("mqtt_temperature_topic", "target/air_temperature", "Air temperature writing topic")
	pflag.Duration("change_interval", 5*time.Second, "How frequently to change the temperature")
	pflag.Float32("home_temperature", 21, "Temperature at home")

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
}

func main() {
	opts := mqtt.NewClientOptions().AddBroker(viper.GetString("mqtt_broker"))
	mqttClient := mqtt.NewClient(opts)

	// Try to connect to the MQTT broker
	token := mqttClient.Connect()
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("Failed to connect to mqtt broker: %v", token.Error())
	}

	// Start the infinite loop
	for {
		// Sleep for the specified interval
		time.Sleep(viper.GetDuration("change_interval"))

		// Generate random temperatures
		temperature := common.MqttTargetAirTemperature{
			Outside: -rand.Float32()*10 - 5, // Random temperature between -15 and -5
			AtHome:  float32(viper.GetFloat64("home_temperature")),
		}

		// Marshal the temperature data into a JSON object
		marshaledTemperature, err := json.Marshal(temperature)
		if err != nil {
			log.Printf("Failed to marshal temperature: %v", err)
			continue
		}

		// Publish the marshaled JSON object to the MQTT topic
		token := mqttClient.Publish(viper.GetString("mqtt_temperature_topic"), 0, false, marshaledTemperature)
		token.Wait()
		if token.Error() != nil {
			log.Printf("Failed to publish message: %v", token.Error())
		}
	}
}
