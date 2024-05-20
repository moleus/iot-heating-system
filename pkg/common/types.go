package common

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

const (
	AirTemperatureTopic   = "target/air_temperature"
	PipesTemperatureTopic = "target/pipes/temperature"
	PipesPressureTopic    = "target/pipes/pressure"
	FuelExpensesTopic     = "target/fuel_expenses"
)
