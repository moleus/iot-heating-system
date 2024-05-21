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
	pflag.String("mqtt_subscribe_temperature_topic", common.AirTemperatureTopic, "Pipes temperature reading topic")
	pflag.String("mqtt_publish_pipes_temperature", common.PipesTemperatureTopic, "Pipes pressure writing topic")
	pflag.String("mqtt_publish_pressure_topic", common.PipesPressureTopic, "Pipes pressure writing topic")

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
	token = mqttClient.Subscribe(viper.GetString("mqtt_subscribe_temperature_topic"), 0, onTemperatureChange)
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("Failed to subscribe to topic: %v", token.Error())
	}

	// Keep the main function running
	select {}
}

func onTemperatureChange(client mqtt.Client, msg mqtt.Message) {
	var airTemperature common.MqttTargetAirTemperature

	err := json.Unmarshal(msg.Payload(), &airTemperature)
	if err != nil {
		log.Printf("Failed to unmarshal message: %v", err)
		return
	}

	inTemperature, outTemperature := calculatePipesTemperature(airTemperature)

	pipesTemperature := common.MqttTargetPipesTemperature{
		InTemperature:  inTemperature,
		OutTemperature: outTemperature,
	}

	marshaledPipesTemperature, err := json.Marshal(pipesTemperature)
	if err != nil {
		log.Printf("Failed to marshal pipes temperature: %v", err)
		return
	}

	token := client.Publish(viper.GetString("mqtt_publish_pipes_temperature"), 0, false, marshaledPipesTemperature)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Failed to publish message: %v", token.Error())
	}

	value := calculatePipePressure(pipesTemperature)

	pressure := common.MqttTargetPipesPressure{
		Value: value,
	}

	marshaledPressure, err := json.Marshal(pressure)
	if err != nil {
		log.Printf("Failed to marshal pressure: %v", err)
		return
	}

	log.Printf("Publishing new pressure value for pipes: %v", pressure)
	token = client.Publish(viper.GetString("mqtt_publish_pressure_topic"), 0, false, marshaledPressure)
	token.Wait()
	if token.Error() != nil {
		log.Printf("Failed to publish message: %v", token.Error())
	}
}

func calculatePipesTemperature(airTemperature common.MqttTargetAirTemperature) (float32, float32) {
	inTemperature := -0.79513742*airTemperature.Outside + 41.788372
	outTemperature := -2.2376321*airTemperature.Outside + 70.973784
	if airTemperature.Outside > common.OutsideTemperatureToDisableSystem {
		return 0, 0
	}
	return inTemperature, outTemperature
}

func calculatePipePressure(temperature common.MqttTargetPipesTemperature) float32 {
	return 1.1 * DesignCapacity * (temperature.OutTemperature - temperature.InTemperature) * SpecificHeatOfWater
}
