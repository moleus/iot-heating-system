package common

type MqttTargetPipesTemperature struct {
	InTemperature  float32 `json:"in_temperature"`
	OutTemperature float32 `json:"out_temperature"`
}

type MqttTargetPipesPressure struct {
	Value float32 `json:"value"`
}
