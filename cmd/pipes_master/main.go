package main

import (
	"encoding/json"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"iot-heating-system/pkg/common"
	"log"
)

// DesignCapacity - МВт
const DesignCapacity = 21.5

// SpecificHeatOfWater - Дж/(кг * C)
const SpecificHeatOfWater = 4190

func init() {
	pflag.String("mqtt_broker", "tcp://mosquitto:1883", "MQTT broker to connect to")
	pflag.String("mqtt_temperature_topic", "target/pipes/temperature", "Pipes temperature reading topic")
	pflag.String("mqtt_pressure_topic", "target/pipes/pressure", "Pipes pressure writing topic")

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

	// Subscribe to the topic
	token = mqttClient.Subscribe(viper.GetString("mqtt_temperature_topic"), 0, onTemperatureChange)
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %v", token.Error())
	}

	// Keep the main function running
	select {}
}

func onTemperatureChange(client mqtt.Client, msg mqtt.Message) {
	var temperature common.MqttTargetPipesTemperature
	err := json.Unmarshal(msg.Payload(), &temperature)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	if temperature.InTemperature >= temperature.OutTemperature {
		log.Printf("Warning: in temperature %f must not be greater to out temperature %f. Skipping", temperature.InTemperature, temperature.OutTemperature)
		return
	}

	value := calculatePipePressure(temperature)

	// Create a Pressure struct and set the value
	pressure := common.MqttTargetPipesPressure{
		Value: value,
	}

	// Marshal the Pressure struct into a JSON object
	marshaledPressure, err := json.Marshal(pressure)
	if err != nil {
		log.Printf("Failed to marshal pressure: %v", err)
		return
	}

	log.Printf("Publishing new pressure value for pipes: %v", pressure)
	// Publish the marshaled JSON object to the output topic
	token := client.Publish(viper.GetString("mqtt_pressure_topic"), 0, false, marshaledPressure)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Failed to publish message: %v", token.Error())
	}
}

// 21.5 -
func calculatePipePressure(temperature common.MqttTargetPipesTemperature) float32 {
	return 1.1 * DesignCapacity * (temperature.OutTemperature - temperature.InTemperature) * SpecificHeatOfWater
}
