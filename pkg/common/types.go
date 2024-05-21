package common

import "time"

type MqttTargetPipesTemperature struct {
	InTemperature  float32 `json:"in_temperature"`
	OutTemperature float32 `json:"out_temperature"`
}

type MqttTargetPipesPressure struct {
	Value float32 `json:"value"`
}

type MqttTargetAirTemperature struct {
	Outside float32 `json:"outside"`
	AtHome  float32 `json:"at_home"`
}

type MqttTargetFuelExpenses struct {
	Value float32 `json:"value"`
}

type MqttTargetFuelExpensesPredictions struct {
	FuelConsumption float32    `json:"temperature,omitempty"`
	Time            *time.Time `json:"time,omitempty"`
}

const (
	AirTemperatureTopic               = "target/air_temperature"
	PipesTemperatureTopic             = "target/pipes/temperature"
	PipesPressureTopic                = "target/pipes/pressure"
	FuelExpensesTopic                 = "target/fuel_expenses"
	WeatherPredictionsTopic           = "predictions/weather"
	FuelExpensesPredictionTopic       = "predictions/fuel_expenses"
	OutsideTemperatureToDisableSystem = 15
)
